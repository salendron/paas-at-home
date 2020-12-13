# SERVICE-ROUTER

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

Service mapping is configured using env vars.
* "SERVICE-NAME": "ADDRESS"

This allows us to deploy multiple copies of this router using different target services
for the same name. That way we could easily build a test environment for certain services,
or do A/B testing and stuff liek that.

## Development
This service is developed using Visual Studio Code and requires the following extensions:
* Docker
* Remote-Containers
* Go

## Deployment
This command runs the service on port 6000 and route two services by setting their target env vars,
TESTHOST_A and TESTHOST_B.
```
docker run -d -p 6000:6000 --name service-router -e PORT='6000' -e TESTHOST_A:'http://localhost:8080' -e TESTHOST_B:'http://localhost:8081' -v /var/run/docker.sock:/var/run/docker.sock --restart unless-stopped service-router:1.0
```

## Usage
To use the router just call this service as you would call the service you need and set header to 
the service name of the target service, as specified in the router'S env vars.
```
curl -i http://localhost:6000/get -H "X-TargetService: TESTHOST_A"
```

curl --header "Content-Type: application/json" --request POST --data '{"definitionPath":"path","name":"theName", "version":"v1"}' http://localhost:6000/register