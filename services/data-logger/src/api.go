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
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

//APIInterface defines the interface of the RESTful API
type APIInterface interface {
	Query(w http.ResponseWriter, r *http.Request)
	Write(w http.ResponseWriter, r *http.Request)
	Collections(w http.ResponseWriter, r *http.Request)
	Initialize(storage StorageInterface)
}

//API implements APIInterface
type API struct {
	Storage StorageInterface
}

//Initialize initializes the API by setting the active storage
func (a *API) Initialize(storage StorageInterface) {
	a.Storage = storage
}

func (a *API) GetDateFilter(name string, r *http.Request) (time.Time, error) {
	val := r.FormValue(name)
	if len(val) == 0 {
		return time.Now().UTC(), errors.New("Value missing")
	}

	return time.Parse(time.RFC3339, val)
}

//API handler to get data items
func (a *API) Query(w http.ResponseWriter, r *http.Request) {
	// Get Request Vars
	vars := mux.Vars(r)
	collectionName, ok := vars["collection"]
	if !ok {
		RaiseError(w, "Collection is missing", http.StatusBadRequest, ErrorCodeCollectionMissing)
		return
	}

	startDate, err := a.GetDateFilter("from", r)
	if err != nil {
		RaiseError(w, "Invalid from-date", http.StatusBadRequest, ErrorCodeInvalidFromDate)
		return
	}

	endDate, err := a.GetDateFilter("to", r)
	if err != nil {
		RaiseError(w, "Invalid to-date", http.StatusBadRequest, ErrorCodeInvalidToDate)
		return
	}

	// Load value
	result, err := a.Storage.ReadData(collectionName, startDate, endDate)
	if err != nil {
		RaiseError(w, "Error loading data.", http.StatusNotFound, ErrorCodeInternal)
		return
	}

	// prepare message
	msg := DataListMessageType{
		Data: result,
	}

	// Write Response
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msg)
}

//API handler to write new data items
func (a *API) Write(w http.ResponseWriter, r *http.Request) {
	// Get Request Vars
	vars := mux.Vars(r)
	collectionName, ok := vars["collection"]
	if !ok {
		RaiseError(w, "Collection is missing", http.StatusBadRequest, ErrorCodeCollectionMissing)
		return
	}

	payload, err := PayloadFromRequestJson(r.Body)
	if err != nil {
		RaiseError(w, "Invalid request body", http.StatusBadRequest, ErrorCodeInvalidRequestBody)
		return
	}

	data, err := a.Storage.WriteData(collectionName, payload)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

//API handler to get collections
func (a *API) Collections(w http.ResponseWriter, r *http.Request) {
	collections, err := a.Storage.ListCollections()
	if err != nil {
		RaiseError(w, fmt.Sprintf("Failed to load collections: %v", err), http.StatusBadRequest, ErrorCodeInternal)
		return
	}

	collectionMessage := CollectionListMessageType{
		Collections: collections,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(collectionMessage)
}
