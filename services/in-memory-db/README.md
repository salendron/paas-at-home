# In-Memory DB
This service implements a very simple in-memory key-value-storage service. It will be used by other services later on, to store data that does not need to exist forever, automatically expirey and needs to be stored and loaded fast.
Yes I know that redis exists, but this whole project is about being as lightweight as possible and also to show how things internally work on a very high level. So let's think about this service as something like a "this is basically how redis works in the most minimal way". We will not implement all functionality of redis, because we do not need it.

This service will be able to:
* Store Data (SET), which will automatically expire.
* Load Data (GET)
* Explicitly delete Data (DELETE)
* List all keys in the DB (LIST-KEYS)

## Development
This service is developed using Visual Studio Code and requires the following extensions:
* Docker
* Remote-Containers
* Go

## API
Description and examples (cUrl) of all API calls and models of this service.

### Models
#### Value
```json
{
        "value":"a value as string",
        "expires-in":180
}
```

#### Error
```json
{
        "error":{
                "message":"No value found for key myrealm/myke",
                "status":404,
                "code":3
                }
}
```

### Methods
#### SET
Sets a value in given realm using given key.

This example sets "a value as string" in realm "myrealm" using key "mykey".
```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"value":"a value as string", "expires-in": 180}' \
  http://localhost:7000/myrealm/mykey
```

#### GET
Gets a value in given realm by given key.

This example get valiue of key "meykey in realm "myrealm".
```
curl -i http://localhost:7000/myrealm/mykey
```