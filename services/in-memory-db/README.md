# In-Memory DB
This service implemnts a very simple in-memory key-value-storage service. It will be used be other services later on to store data that does not need to exist forever, automatically expire and needs to be stored and loaded fast.
Yes I know that redis exists, but this whole project is about beeing as lightweight as possible and also to show how things internally work on a very high level. So let's think about this service as something like a "this is basically how redis works in the most minimal way". We will not implement all functionality of redis, because we do not need it.

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

*TODO WRITE SERVICE DOCUMENTATION AS SOON AS SERVICE IS IMPLEMENTED*