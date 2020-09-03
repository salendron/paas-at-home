# IN-MEMORY-DB
in-memory-db is a service that does something like redis on a very basic
level.
It implements a very basic key/value storage that can be used to store
data that does not need to be persistet, because the service doesn't do that,
but has to be saved and loaded fast. It is also implemented to automatically
delete data based on an expiration time. There is no way to store data permanently!
Data is lost either after the service restarts or after the set expiration time
is over.
To structure data bit better it implements realms, which is just one layer more
to devide data into seperate spaces. This can be used to seperate storage spaces
for services using this, to eliminate the problem of key conflicts.
You can use this if you want a very lightweight in-memory key/value storage
and redis is just too much, or use it to see how key/value databases could be
implemented in a very basic way. It also shows how you can use go routines to do
things after a set amount of time asynchronously.

This service is be able to:
* Store Data (SET), which will automatically expire.
* Load Data (GET)
* Explicitly delete Data (DELETE)
* List all keys in a realm (LIST-KEYS)
* List all realms (LIST-REALMS)

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

#### Key List
```json
{
        "keys":["key1", "key2", ...]
}
```

#### Realm List
```json
{
        "realms":["realm1", "realm2", ...]
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

This example gets the value of key "meykey" in realm "myrealm".
```
curl -i http://localhost:7000/myrealm/mykey
```

#### DELETE
Deletes a value in given realm using given key.

This example deletes the value in realm "myrealm" with the key "mykey".
```
curl --request DELETE http://localhost:7000/myrealm/mykey
```

#### GET Keys 
Gets all keys in a given realm.

This example gets all keys in realm "myrealm".
```
curl -i http://localhost:7000/myrealm/keys   
```

#### GET Realms 
Gets all realms.

This example gets all realms.
```
curl -i http://localhost:7000/realms  
```
