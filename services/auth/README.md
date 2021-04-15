# AUTH
Auth is a service that can be used to authenticate user and retrieve
permissions. Login returns an access and a refresh token. The access token
can be used to autheticate this user in other services, which can call decode 
to verify the token. Decode retirns the user's id as well as all permissions
of this user. The refresh token can be used to refresh a user'S authentication.
Again, do not use this in production, but it is a nice example on how to
implment JWT authentication and a refresh mechanism.

## Deployment
This command runs the service on port 7004 and mounts the local directory /media/external/storage/auth to /data which will be used by the service to read user data from. IT also sets all secrets and also the password for the super user (su).
```
docker run -d -p 7004:7004 --name auth -e PORT='7004' -e DATA_DIRECTORY='/data' -e AUTH_ACCESS_SECRET='myGreatAccessSecret' -e AUTH_REFRESH_SECRET='myGreatRefreshSecret' -e AUTH_SERVICE_SECRET='myGreatServiceSecret' -e SU_PWD='myGreatSuPassword' -v /var/run/docker.sock:/var/run/docker.sock --restart unless-stopped --mount type=bind,source=/media/external/storage/auth,target=/data data-logger:1.0
```

## API
Description and examples (cUrl) of all API calls and models of this service

### Errors
All errors are served as an object like this and will return a suitable HTTP
status code.
```json
{
        "error":{
                "message":"Invalid Token",
                "status":403,
                "code":3
                }
}
```

### Methods
#### LOGIN
Logs in a user using username and password.
```
curl --header "Content-Type: application/json" \
        --request POST \ 
        --data '{"username":"theUsername","password":"thePAssword"}' \
        http://localhost:7004/login
```

This call returns a access-token and a refresh token for this session.
Example Response:
```json
{
        "access-token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOiIyMDIwLTEwLTI1VDE1OjI4OjA2Ljk5MjE3MTA4MVoiLCJwZXJtaXNzaW9ucyI6Ilt7XCJrZXlcIjpcIlJPT1RcIixcIm1ldGFcIjpudWxsfV0iLCJ1c2VyX2lkIjoic3UifQ.GB9skCdYbnbDsA9IcMiybMRmgNlt4P_F-inUfuvaXZk","refresh-token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOiIyMDIwLTExLTAxVDE1OjEzOjA2Ljk5MjE3NDA2M1oiLCJ1c2VyX2lkIjoic3UifQ.ugCGoLxSr2XoJxCTedRVM1mFsT-LZgRh-uyEwTraOsQ"
}
```

#### REFRESH TOKEN
This call generates a new access and refresh token for the session with given refresh token.
```
curl --header "Content-Type: application/json" \
        --request POST \
        --data '{"refresh-token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOiIyMDIwLTExLTAxVDE1OjEzOjA2Ljk5MjE3NDA2M1oiLCJ1c2VyX2lkIjoic3UifQ.ugCGoLxSr2XoJxCTedRVM1mFsT-LZgRh-uyEwTraOsQ"}' \ 
        http://localhost:7004/refresh
```
Example Response:
```json
{
        "access-token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOiIyMDIwLTEwLTI1VDE1OjMyOjIxLjMzODgyMjIxNloiLCJwZXJtaXNzaW9ucyI6Ilt7XCJrZXlcIjpcIlJPT1RcIixcIm1ldGFcIjpudWxsfV0iLCJ1c2VyX2lkIjoic3UifQ.N7dloh8YAvyBvEck36Q7moMH1MWNU8iW11A3xNKhLto","refresh-token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOiIyMDIwLTExLTAxVDE1OjE3OjIxLjMzODgyMzkzOFoiLCJ1c2VyX2lkIjoic3UifQ.03r_7ImXxLGElpRVpgHCO4jtErKl63SJaF7CEf0yka8"
}
```

#### DECODE TOKEN
This call is used to verify and decode a given access token.

```
curl --header "Content-Type: application/json" \
        --request POST \ 
        --data '{"access-token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOiIyMDIwLTEwLTI1VDE1OjMyOjIxLjMzODgyMjIxNloiLCJwZXJtaXNzaW9ucyI6Ilt7XCJrZXlcIjpcIlJPT1RcIixcIm1ldGFcIjpudWxsfV0iLCJ1c2VyX2lkIjoic3UifQ.N7dloh8YAvyBvEck36Q7moMH1MWNU8iW11A3xNKhLto"}' \
        http://localhost:7004/decode
```
Example Response:
```json
{
        "user-id":"aUserId",
        "permissions":[
                {"key":"ROOT","meta":null}
        ],
        "expires":"2020-10-25T15:32:21.338822216Z"
}
```