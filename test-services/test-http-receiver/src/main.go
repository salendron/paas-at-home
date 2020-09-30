/*
TEST-HTTP-RECEIVER

test-http-receiver is a service to test http requests. It handles GET, POST, PUT
and DELETE requests and simply logs all information about the request incl. request
body.

###################################################################################

main.go
This is the main entrypoint of the service. It starts the service and
routes all handlers.

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
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

func logRequest(w http.ResponseWriter, r *http.Request) {
	logMsg := make([]string, 0)

	// Request
	method := r.Method
	logMsg = append(logMsg, fmt.Sprintf("%v: %v", r.Method, r.URL))
	logMsg = append(logMsg, fmt.Sprintf("Remote Address: %v", r.RemoteAddr))
	logMsg = append(logMsg, fmt.Sprintf("ContentLength: %v", r.ContentLength))

	// Headers
	logMsg = append(logMsg, "Headers:")
	for name, values := range r.Header {
		// Loop over all values for the name.
		for _, value := range values {
			logMsg = append(logMsg, fmt.Sprintf("\t%v: %v", name, value))
		}
	}

	// Body
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err == nil {
		logMsg = append(logMsg, fmt.Sprintf("Request Body: %v", buf.String()))
	}

	log.Println(strings.Join(logMsg, "\n"))

	// Respond
	w.Header().Add("Content-Type", "application/json")

	switch {
	case method == "GET":
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(strings.Join(logMsg, "\n"))
	case method == "POST":
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(strings.Join(logMsg, "\n"))
	case method == "PUT":
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(strings.Join(logMsg, "\n"))
	case method == "DELETE":
		w.WriteHeader(http.StatusNoContent)
	}
}

//main is the main entrypoint of the service.
func main() {
	r := mux.NewRouter()

	// Topic management
	r.HandleFunc("/get", logRequest).Methods("GET")
	r.HandleFunc("/post", logRequest).Methods("POST")
	r.HandleFunc("/put", logRequest).Methods("PUT")
	r.HandleFunc("/delete", logRequest).Methods("DELETE")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", os.Getenv("PORT")), r))
}
