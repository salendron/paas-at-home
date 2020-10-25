/*
api.go
Defines the RESTful API interface of this application and implements all api
methods.

###################################################################################

MIT License

Copyright (c) 2020 Bruno Hautzenberger

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// APIInterface defines the interface of the RESTful API
type APIInterface interface {
	Initialize(storage StorageInterface, tokenbuilder TokenBuilderInterface)
	UserLogin(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
	DecodeToken(w http.ResponseWriter, r *http.Request)
	ServiceLogin(w http.ResponseWriter, r *http.Request)
}

// API implements APIInterface
type API struct {
	Storage      StorageInterface
	Tokenbuilder TokenBuilderInterface
}

// Initialize initializes the API by setting the active storage and publisher
func (a *API) Initialize(storage StorageInterface, tokenbuilder TokenBuilderInterface) {
	a.Storage = storage
	a.Tokenbuilder = tokenbuilder
}

// parseRequestPayload parses the given json data of the request's io.ReadCloser
func parseRequestPayload(rc io.ReadCloser, dst interface{}) error {
	err := json.NewDecoder(rc).Decode(dst)
	if err != nil {
		return err
	}

	return nil
}

// UserLogin handles user login api requests
func (a *API) UserLogin(w http.ResponseWriter, r *http.Request) {
	loginMsg := &UserLoginType{}
	err := parseRequestPayload(r.Body, loginMsg)
	if err != nil {
		RaiseError(w, "Invalid request body. Invalid Json format", http.StatusBadRequest, ErrorCodeInvalidRequestBody)
		return
	}

	user, ok, err := a.Storage.GetUserByCredentials(loginMsg.Username, loginMsg.Password)
	if err != nil {
		RaiseError(w, err.Error(), http.StatusInternalServerError, ErrorCodeInternal)
		return
	}

	if !ok {
		RaiseError(w, "Login failed", http.StatusUnauthorized, ErrorCodeLoginFailed)
		return
	}

	td, err := a.Tokenbuilder.CreateUserToken(user)
	if err != nil {
		RaiseError(w, err.Error(), http.StatusInternalServerError, ErrorCodeInternal)
		return
	}

	resp := UserTokenType{
		AccessToken:  td.AccessToken,
		RefreshToken: td.RefreshToken,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// RefreshToken is the API Handler for token refresh requests
func (a *API) RefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshMsg := &RefreshTokenRequestType{}
	err := parseRequestPayload(r.Body, refreshMsg)
	if err != nil {
		RaiseError(w, "Invalid request body. Invalid Json format", http.StatusBadRequest, ErrorCodeInternal)
		return
	}

	token, err := jwt.Parse(refreshMsg.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v", token.Header["alg"])
		}
		return []byte(os.Getenv("AUTH_REFRESH_SECRET")), nil
	})
	if err != nil {
		RaiseError(w, err.Error(), http.StatusBadRequest, ErrorCodeUnexpectedSigningMethod)
		return
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		RaiseError(w, "Invalid Token", http.StatusBadRequest, ErrorCodeInvalidToken)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		RaiseError(w, "Invalid token claims", http.StatusBadRequest, ErrorCodeInvalidToken)
		return
	}

	exp, ok := claims["exp"].(time.Time)
	if !ok {
		RaiseError(w, "Missing exp", http.StatusBadRequest, ErrorCodeInvalidToken)
		return
	}

	if exp.After(time.Now().UTC()) {
		RaiseError(w, "Token expired", http.StatusUnauthorized, ErrorCodeTokenExpired)
		return
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		RaiseError(w, "Missing user_id", http.StatusBadRequest, ErrorCodeInvalidToken)
		return
	}

	user, err := a.Storage.GetUser(userID)
	if err != nil {
		RaiseError(w, err.Error(), http.StatusInternalServerError, ErrorCodeInternal)
		return
	}

	if user == nil {
		RaiseError(w, fmt.Sprintf("Unknown user %v", userID), http.StatusBadRequest, ErrorCodeInvalidToken)
		return
	}

	td, err := a.Tokenbuilder.CreateUserToken(user)
	if err != nil {
		RaiseError(w, err.Error(), http.StatusInternalServerError, ErrorCodeInternal)
		return
	}

	resp := &UserTokenType{
		AccessToken:  td.AccessToken,
		RefreshToken: td.RefreshToken,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// DecodeToken is the API handler to decode and verify access tokens
func (a *API) DecodeToken(w http.ResponseWriter, r *http.Request) {
	decodeMsg := &DecodeTokenMessage{}
	err := parseRequestPayload(r.Body, decodeMsg)
	if err != nil {
		RaiseError(w, "Invalid request body. Invalid json format", http.StatusBadRequest, ErrorCodeInvalidRequestBody)
		return
	}

	token, err := jwt.Parse(decodeMsg.AccessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v", token.Header["alg"])
		}
		return []byte(os.Getenv("AUTH_ACCESS_SECRET")), nil
	})
	if err != nil {
		RaiseError(w, err.Error(), http.StatusBadRequest, ErrorCodeUnexpectedSigningMethod)
		return
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		RaiseError(w, "Invalid token", http.StatusBadRequest, ErrorCodeInvalidToken)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		RaiseError(w, "Invalid token claims", http.StatusBadRequest, ErrorCodeInvalidToken)
		return
	}

	expValue, ok := claims["exp"].(string)
	if !ok {
		RaiseError(w, "Missing exp", http.StatusBadRequest, ErrorCodeInvalidToken)
		return
	}

	exp, _ := time.Parse(time.RFC3339, expValue)
	if exp.Before(time.Now().UTC()) {
		RaiseError(w, "Token expired", http.StatusBadRequest, ErrorCodeTokenExpired)
		return
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		RaiseError(w, "Missing user id", http.StatusBadRequest, ErrorCodeInvalidToken)
		return
	}

	permissionsValue, ok := claims["permissions"].(string)
	permissions := make([]Permission, 0)
	err = json.Unmarshal([]byte(permissionsValue), &permissions)
	if err != nil {
		RaiseError(w, "Missing permissions", http.StatusBadRequest, ErrorCodeInvalidToken)
		return
	}

	decodedToken := DecodedTokenMessage{
		UserID:      userID,
		Permissions: permissions,
		Expires:     exp,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(decodedToken)
}

// ServiceLogin handles service login api requests
func (a *API) ServiceLogin(w http.ResponseWriter, r *http.Request) {
	loginMsg := ServiceLoginType{}
	err := parseRequestPayload(r.Body, loginMsg)
	if err != nil {
		RaiseError(w, "Invalid request body. Invalid json format", http.StatusBadRequest, ErrorCodeInvalidRequestBody)
		return
	}
	service, ok, err := a.Storage.GetServiceByCredentials(loginMsg.ID, loginMsg.Key)
	if err != nil {
		RaiseError(w, err.Error(), http.StatusInternalServerError, ErrorCodeInternal)
		return
	}

	if !ok {
		RaiseError(w, "Login failed", http.StatusUnauthorized, ErrorCodeLoginFailed)
		return
	}

	td, err := a.Tokenbuilder.CreateServiceToken(service)
	if err != nil {
		RaiseError(w, err.Error(), http.StatusInternalServerError, ErrorCodeInternal)
		return
	}

	resp := &ServiceTokenType{
		AccessToken: td.AccessToken,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
