Tea Selector API
===

This API interacts with a database, containing information used for the tea selector app.

## Dependencies
- [gorilla/mux](https://github.com/gorilla/mux) - `go get -u github.com/gorilla/mux`

## Building

To build an executable, run:
```
go build
```

## Configuration
The `config.yml` file gives an example configuration. This can be changed to your liking. You can:
- Set the port the API runs on.
- Set the database location
- Set the default tea types and owners.

## Interacting with the API
By default, the API will be running on `localhost:7344`.

### Tea Types

- To see all current tea types, send a GET request to `/types`.
- To get information about a tea type, send a GET request: `/type/{id}` 

- To add a new tea type, send a POST request to `/type`. An example body is :

        {
            "name": "Black Tea"
        }

- To delete a tea type, send a DELETE request: `/type/{id}`
