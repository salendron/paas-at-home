/*
AUTH

Auth is a service that can be used to authenticate user and retrieve
permissions. Login returns an access and a refresh token. The access token
can be used to autheticate this user in other services, which can call decode
to verify the token. Decode retirns the user's id as well as all permissions
of this user. The refresh token can be used to refresh a user'S authentication.
Again, do not use this in production, but it is a nice example on how to
implment JWT authentication and a refresh mechanism.

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
	r.HandleFunc("/refresh", api.RefreshToken).Methods("POST")
	r.HandleFunc("/servicelogin", api.ServiceLogin).Methods("POST")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", os.Getenv("PORT")), r))
}
