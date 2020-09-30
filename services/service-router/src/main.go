/*
SERVICE-ROUTER

service-router is a service that can be used to route requests based on X-TargetService
header.
It is basically a reverse proxy used to route requests to a service without the need to
know about the service's address or port. In this context it is used for services to call
each other by name (X-TargetService header).
This allows us to quickly swap services, redeploy them somewhere else and stuff like that,
without the need of having to reconfigure all other services that rely on them.
Services always call this service with X-TargetService header set, to request the service
they actually need and this service will relay the request and return the response.
Path, request body and headers will be forwarded as well.

###################################################################################

main.go
This is the main entrypoint of the service. It starts the service and
starts the reverse proxy.

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
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

// getTargetUrl tries to get the target URL specified for given service name.
// returns an error if there is no service configured for given target name.
func getTargetUrl(targetName string) (string, error) {
	cleanTargetName := strings.ToUpper(strings.TrimSpace(targetName))

	if targetUrl, ok := os.LookupEnv(cleanTargetName); ok {
		return targetUrl, nil
	}

	return "", errors.New(fmt.Sprintf("No target specified for '%v'", cleanTargetName))
}

// serveProxy starts the - non-blocking - reverse proxy for the current request
func serveProxy(targetUrl string, w http.ResponseWriter, r *http.Request) error {
	url, err := url.Parse(targetUrl)
	if err != nil {
		return err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = url.Host

	//non blocking
	proxy.ServeHTTP(w, r)

	return nil
}

// handleRedirect handles all incoming requests. It reads X-TargetService
// header to determine which target service to use and then proceeds to start
// a proxy using serveProxy function.
func handleRedirect(w http.ResponseWriter, r *http.Request) {
	targetName := r.Header.Get("X-TargetService")
	if len(targetName) == 0 {
		RaiseError(w, "Missing X-TargetService", http.StatusBadRequest, ErrorCodeMissingTargetServiceName)
		return
	}

	targetUrl, err := getTargetUrl(targetName)
	if err != nil {
		RaiseError(w, fmt.Sprintf("No service specified for target %v", targetName), http.StatusBadRequest, ErrorCodeUnknownTargetService)
		return
	}

	serveProxy(targetUrl, w, r)
}

// main is the main entrypoint of the service. It starts the server on PORT
// specified in env vars and routes everthing to handleRedirect.
func main() {
	// Serve router
	http.HandleFunc("/", handleRedirect)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", os.Getenv("PORT")), nil))
}
