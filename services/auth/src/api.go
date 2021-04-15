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
	"github.com/labstack/echo"
)

// APIInterface defines the interface of the RESTful API
type APIInterface interface {
	Initialize(storage StorageInterface, tokenbuilder TokenBuilderInterface)
	UserLogin(c echo.Context) error
	RefreshToken(c echo.Context) error
	DecodeToken(c echo.Context) error
}

// API implements APIInterface
type API struct {
	Storage      StorageInterface
	Tokenbuilder TokenBuilderInterface
}

// Initialize initializes the API by setting the active storage and tokenbuilder
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
func (a *API) UserLogin(c echo.Context) error {
	loginMsg := new(UserLoginType)
	err := c.Bind(loginMsg)
	if err != nil {
		return err
	}

	user, err := a.Storage.GetUserByCredentials(loginMsg.Username, loginMsg.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Login failed")
	}

	td, err := a.Tokenbuilder.CreateUserToken(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	resp := UserTokenType{
		AccessToken:  td.AccessToken,
		RefreshToken: td.RefreshToken,
	}

	return c.JSON(http.StatusOK, resp)
}

// RefreshToken is the API Handler for token refresh requests
func (a *API) RefreshToken(c echo.Context) error {
	refreshMsg := new(RefreshTokenRequestType)
	err := c.Bind(refreshMsg)
	if err != nil {
		return err
	}

	token, err := jwt.Parse(refreshMsg.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v", token.Header["alg"])
		}
		return []byte(os.Getenv("AUTH_REFRESH_SECRET")), nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Token Claims")
	}

	expValue, ok := claims["exp"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing exp")
	}

	exp, _ := time.Parse(time.RFC3339, expValue)
	if exp.Before(time.Now().UTC()) {
		return echo.NewHTTPError(http.StatusUnauthorized, "Token expired")
	}

	userID, ok := claims["user_id"].(int)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing user_id")
	}

	user, err := a.Storage.GetUserByID(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if user == nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unknown user %v", userID))
	}

	td, err := a.Tokenbuilder.CreateUserToken(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	resp := &UserTokenType{
		AccessToken:  td.AccessToken,
		RefreshToken: td.RefreshToken,
	}

	return c.JSON(http.StatusOK, resp)
}

// DecodeToken is the API handler to decode and verify access tokens
func (a *API) DecodeToken(c echo.Context) error {
	decodeMsg := new(DecodeTokenMessage)
	err := c.Bind(decodeMsg)
	if err != nil {
		return err
	}

	token, err := jwt.Parse(decodeMsg.AccessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v", token.Header["alg"])
		}
		return []byte(os.Getenv("AUTH_ACCESS_SECRET")), nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Token Claims")
	}

	expValue, ok := claims["exp"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing exp")
	}

	exp, _ := time.Parse(time.RFC3339, expValue)
	if exp.Before(time.Now().UTC()) {
		return echo.NewHTTPError(http.StatusBadRequest, "Token expired")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing user_id")
	}

	permissionsValue, ok := claims["permissions"].(string)
	permissions := make([]Permission, 0)
	err = json.Unmarshal([]byte(permissionsValue), &permissions)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing permissions")
	}

	decodedToken := DecodedTokenMessage{
		UserID:      userID,
		Permissions: permissions,
		Expires:     exp,
	}

	return c.JSON(http.StatusOK, decodedToken)
}
