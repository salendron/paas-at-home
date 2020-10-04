# DATA-LOGGER
data-logger is a service that can be used to store custom json data items
such as logs or sensor data or what ever you want.
To structure data by type this service implments collections that are dynamically
created as soon as data is saved to a collection, specified by name.
It is important to know that data can only queried by collection and timeframe.
Also data can't be deleted, so consider this as a long term storage for immutable
data.
You can use this if you want a very lightweight json data storage for your services.
It shows how you can split data into seperate data files, read query params using
mux and also how to lock files during writes using sync.Mutex.

This service is be able to:
* Store Data 
* Query Data
* List all collections

## Development
This service is developed using Visual Studio Code and requires the following extensions:
* Docker
* Remote-Containers
* Go

## Deployment
This command runs the service on port 7001 and mounts the local directory /media/external/storage/data-logger to /data
which will be used by the service to write the data files to.
```
docker run -d -p 7001:7001 --name data-logger -e PORT='7001' -e DATA_DIRECTORY='/data' -v /var/run/docker.sock:/var/run/docker.sock --restart unless-stopped --mount type=bind,source=/media/external/storage/data-logger,target=/data data-logger:1.0
```

## API
Description and examples (cUrl) of all API calls and models of this service.

### Models
#### Date Item Wrapper Type
All logged data items get wrapped into a uniform structure that contains a UUID for this item and
also a created-at timestamp. The original raw data item can be found in "payload".
```json
{
        "uuid":"2020-09-03-b8e33721-eea9-41fc-8173-98c0b23b5dad",
        "created-at":"2020-09-03T18:07:10.427571Z",
        "payload":{"valueA":"some custom value","valueB":42}
}
```

#### Query Result
```json
{
        "data":[
                {Data-Item-Wrapper-Item 1}, {Data-Item-Wrapper-Item 2}, ...
        ]
}
```

#### Collection List
```json
{
        "collections":["collection1", "collection2", ...]
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
#### WRITE
Writes a new data item to a collection. Not existing collections will automatically be created.

This example creates a new data item in collection "mycollection". It will return the created item wrapped in
the data item wrapper structure.
```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"valueA":"some custom value", "valueB": 42}' \
  http://localhost:7001/mycollection
```

#### QUERY
Query for data items in a collection in a time range.
u
This example gets all data items between 2020-09-01T10:30:00Z and 2020-09-03T22:45:00Z in collection "testCollection".
```
curl -i 'http://localhost:7001/testCollection?from=2020-09-01T10:30:00Z&to=2020-09-03T22:45:00Z'
```

#### GET COLLECTIONS
Gets all collections.

```
curl -i http://localhost:7001/info/collections
 
```