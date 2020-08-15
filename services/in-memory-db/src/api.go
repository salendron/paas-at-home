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
	"net/http"

	"github.com/gorilla/mux"
)

type APIInterface interface {
	Get(w http.ResponseWriter, r *http.Request)
	Set(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Keys(w http.ResponseWriter, r *http.Request)
	Realms(w http.ResponseWriter, r *http.Request)
	Initialize(storage StorageInterface)
}

type API struct {
	Storage StorageInterface
}

func (a *API) Initialize(storage StorageInterface) {
	a.Storage = storage
}

func (a *API) Get(w http.ResponseWriter, r *http.Request) {
	// Get Request Vars
	vars := mux.Vars(r)
	realm, ok := vars["realm"]
	if !ok {
		RaiseError(w, "Realm is missing", http.StatusBadRequest, ErrorCodeRealmMissing)
		return
	}

	key, ok := vars["key"]
	if !ok {
		RaiseError(w, "Key is missing", http.StatusBadRequest, ErrorCodeKeyMissing)
		return
	}

	// Load value
	ok, value := a.Storage.Get(realm, key)
	if !ok {
		RaiseError(w, fmt.Sprintf("No value found for key %v/%v", realm, key), http.StatusNotFound, ErrorCodeEntityNotFound)
		return
	}

	// Write Response
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(value.ToValueMessageType())
}

func (a *API) Set(w http.ResponseWriter, r *http.Request) {
	// Get Request Vars
	vars := mux.Vars(r)
	realm, ok := vars["realm"]
	if !ok {
		RaiseError(w, "Realm is missing", http.StatusBadRequest, ErrorCodeRealmMissing)
		return
	}

	key, ok := vars["key"]
	if !ok {
		RaiseError(w, "Key is missing", http.StatusBadRequest, ErrorCodeKeyMissing)
		return
	}

	value, err := ValueFromValueMessageType(r.Body)
	if err != nil {
		RaiseError(w, "Invalid request body", http.StatusBadRequest, ErrorCodeInvalidRequestBody)
		return
	}

	a.Storage.Set(realm, key, value)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(value.ToValueMessageType())
}

func (a *API) Delete(w http.ResponseWriter, r *http.Request) {
	// Get Request Vars
	vars := mux.Vars(r)
	realm, ok := vars["realm"]
	if !ok {
		RaiseError(w, "Realm is missing", http.StatusBadRequest, ErrorCodeRealmMissing)
		return
	}

	key, ok := vars["key"]
	if !ok {
		RaiseError(w, "Key is missing", http.StatusBadRequest, ErrorCodeKeyMissing)
		return
	}

	ok = a.Storage.Delete(realm, key)
	if !ok {
		RaiseError(w, fmt.Sprintf("No value found for key %v/%v", realm, key), http.StatusNotFound, ErrorCodeEntityNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) Keys(w http.ResponseWriter, r *http.Request) {
	// Get Request Vars
	vars := mux.Vars(r)
	realm, ok := vars["realm"]
	if !ok {
		RaiseError(w, "Realm is missing", http.StatusBadRequest, ErrorCodeRealmMissing)
		return
	}

	keys := a.Storage.Keys(realm)
	keysMessage := KeyListMessageType{
		Keys: keys,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(keysMessage)
}

func (a *API) Realms(w http.ResponseWriter, r *http.Request) {
	realms := a.Storage.Realms()
	keysMessage := RealmListMessageType{
		Realms: realms,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(keysMessage)
}
