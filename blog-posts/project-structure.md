# Project structure

The structure of the project is going to be vary depending on how complex the service you are building will be. In this project we have tried to keep the structure quite flat to ensure the best possible readability.

### /cmd
This is where we put our main applications. The directories inside should have the same name as the executable generated.
This will contain your main application, but could also contain additional tooling like a migration application.
The `main.go` should be kept as small as possible, only calling the methods to bootstrap the service.

### /api
This stores the router and handlers for the REST api.

### /migrations
This contains the database migration files.
### /storage

## Interesting links
https://www.youtube.com/watch?v=oL6JBUk6tj0 // https://github.com/katzien/go-structure-examples/blob/master/modular/storage/storage.go
