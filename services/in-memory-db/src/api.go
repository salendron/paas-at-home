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
	response := ValueMessageType{
		Value: value, //TODO expires in
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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

	msg := ValueMessageType{}
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		RaiseError(w, "Key is missing", http.StatusBadRequest, ErrorCodeInvalidRequestBody)
	}

	a.Storage.Set(realm, key, msg.Value, msg.ExpiresIn)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msg)
}

func (a *API) Delete(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (a *API) Keys(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (a *API) Realms(w http.ResponseWriter, r *http.Request) {
	// TODO
}
