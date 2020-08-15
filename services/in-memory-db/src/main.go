/*
IN-MEMORY-DB

in-memory-db is a service that does something like redis on a very basic
level.
It implements a very basic key/value storage that can be used to store
data that does not need to be persistet, because the service doesn't od that,
but has to be saved and loaded fast. It is also implemented to automatically
delete data based on an expiration time. There is no way to store data permanently!
Data is lost either after the service restarts or after the set expiration time
is over.
To structure data bit better it implements realms, which is just one layer more
to devide data into seperate spaces. This can be used to seperate storage spaces
for services using this, to eliminate the problem of of key conflicts.
You can use this if you want a very lightweight in-memory key/value storage
and redis is just too much, or use it to see how key/value databases could be
implemented in a very basic way. It also shows how you can use go routines to do
things after a set amount of time asynchronously.

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
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var storage StorageInterface = &Storage{}
var api *API = &API{}

//init initializes storage and api
func init() {
	storage.Initialize()
	api.Initialize(storage)
}

//main is the main entrypoint of the service. It routes all API methods
//and starts the server on PORT specified in env vars.
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
