package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var storage StorageInterface = &Storage{}
var api *API = &API{}

func init() {
	storage.Initialize()
	api.Initialize(storage)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/{realm}/{key}", api.Get).Methods("GET")
	r.HandleFunc("/{realm}/{key}", api.Set).Methods("POST")
	r.HandleFunc("/{realm}/{key}", api.Delete).Methods("DELETE")
	r.HandleFunc("/{realm}/keys", api.Keys).Methods("GET")
	r.HandleFunc("/realms", api.Realms).Methods("GET")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", os.Getenv("PORT")), r))
}
