/*
AUTH

Auth is a service that can be used to authenticate user and retrieve
permissions. Login returns an access and a refresh token. The access token
can be used to autheticate this user in other services, which can call decode
to verify the token. Decode returns the user's id as well as all permissions
of this user. The refresh token can be used to refresh a user's authentication.
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
	"os"

	"github.com/labstack/echo"
)

var storage StorageInterface = nil
var tokenbuilder TokenBuilderInterface = &TokenBuilder{}
var api APIInterface = &API{}

//main is the main entrypoint of the service. It routes all API methods
//and starts the server on PORT specified in env vars.
func main() {
	dbStorage := &Storage{}
	err := dbStorage.Connect(os.Getenv("DSN"))
	if err != nil {
		log.Fatalf("DB Connection failed: %v", err)
	}
	storage = dbStorage

	api.Initialize(storage, tokenbuilder)

	e := echo.New()

	e.POST("/login", api.UserLogin)
	e.POST("/decode", api.DecodeToken)
	e.POST("/refresh", api.RefreshToken)

	// Bind to a port and pass our router in
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", os.Getenv("PORT"))))
}
