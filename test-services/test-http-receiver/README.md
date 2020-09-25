# TEST-HTTP-RECEIVER
test-http-receiver is a service to test http requests. It handles GET, POST, PUT
and DELETE requests and simply logs all information about the request incl. request
body.

## Development
This service is developed using Visual Studio Code and requires the following extensions:
* Docker
* Remote-Containers
* Go

## Deployment
This command runs the service on port 8080.
```
docker run -d -p 8080:8080 --name test-http-receiver -e PORT='8080' -v /var/run/docker.sock:/var/run/docker.sock --restart unless-stopped test-http-receiver:1.0
```

### Methods
#### GET 
A TestEndpoint for GET requests.

```
curl -i http://localhost:8080/get -H "X-TestHeader: 123"
 
```

#### POST
A TestEndpoint for POST requests.

```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"valueA":"some custom value", "valueB": 42}' \
  http://localhost:8080/post
```

#### PUT
A TestEndpoint for PUT requests.

```
curl --header "Content-Type: application/json" \
  --request PUT \
  --data '{"valueA":"some custom value", "valueB": 42}' \
  http://localhost:8080/put
```

#### DELETE
A TestEndpoint for DELETE requests.

```
curl --request DELETE http://localhost:8080/delete
```