/*
DATA-LOGGER

data-logger is a service that can be used to store custom json data items
such as logs or sensor data or what ever you want.
To structure data by type this service implments collections that are dynamically
created as soon as data is saved to a collection, specified by name.
It is important to know that data can only queried by collection and timeframe.
Also data can't be deleted, so consider this as a long term storage for immutable
data.
You can use this if you want a very lightweight json data storage for your services.
It shows how you can split data into seperate data files, read query params using
mux and also how to lock files during writes using sync.Mutex.

###################################################################################

main.go
This is the main entrypoint of the service. It starts the service and
routes all API methods.

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
	"fmt"
	"os"
	"net/http"
	"github.com/gorilla/mux"
	"log"
)

var storage StorageInterface = &Storage{}
var tokenbuilder TokenBuilderInterface = &TokenBuilder{}
var api APIInterface = &API{}

//init initializes storage and api
func init() {
	storage.Initialize(os.Getenv("DATA_DIRECTORY"))
	api.Initialize(storage, tokenbuilder)
}

//main is the main entrypoint of the service. It routes all API methods
//and starts the server on PORT specified in env vars.
func main() {
	r := mux.NewRouter()

	r.HandleFunc("/login", api.UserLogin).Methods("POST")
	r.HandleFunc("/decode", api.DecodeToken).Methods("POST")
	r.HandleFunc("/servicelogin", api.ServiceLogin).Methods("POST")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", os.Getenv("PORT")), r))
}
