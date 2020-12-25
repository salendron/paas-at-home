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

 ./oapi-codegen /home/ubuntu/git/paas-at-home/services/service-router/openapi.json > /home/ubuntu/git/paas-at-home/services/service-router/src/servicerouter.gen.go

*/

package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	var storage CachedSQLiteStorage
	storage.Initialize("./test.db")

	var api API
	api.SetStorage(&storage)

	e := echo.New()
	RegisterHandlers(e, &api)

	// Start HealthChecks
	healthCheck := &HealthCheck{}
	healthCheck.DoHealthChecks(&storage)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", os.Getenv("PORT"))))
}
