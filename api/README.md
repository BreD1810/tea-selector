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

Additionally, `tea-store.sql` is included to setup an example database. To use it, run `sqlite3 tea-store.db`, and then `.read tea-store.sql`.

## Interacting with the API
By default, the API will be running on `localhost:7344`.

### Tea Types
- To see all current tea types, send a GET request to `/types`
- To see all teas of all types, send a GET request to `/types/teas`
- To get information about a tea type, send a GET request: `/type/{id}` 
- To add a new tea type, send a POST request to `/type`. An example body is:

        {
            "name": "Black Tea"
        }

- To delete a tea type, send a DELETE request: `/type/{id}`

### Owners
- To see all current owners, send a GET request to `/owners`
- To get information about an owner, send a GET request: `/owner/{id}`
- To add a new owner, send a POST request to `/owner`. An example body is:

        {
            "name": "John"
        }

- To delete an owner, send a DELETE request: `/owner/{id}`

### Tea
- To see all teas, send a GET request to `/teas`
- To get information about a tea, send a GET request: `/tea/{id}`
- To add a new tea, send a POST request to `/tea`. An example body is:

        {
            "name": "Snowball",
            "type": {
                "id": 1
            }
        }

- To delete a tea, send a DELETE request to `/tea/{id}`

### Tea Owners
- To see all teas with all their owners, send a GET request to `/teas/owners/`
- To see the owners for a specific tea, send a GET request: `/tea/{id}/owners`
- To add an owner, send a POST request to `/tea/{id}/owner`. An example body is:

        {
            "id": 1
        }

- Tp delete an owner from a tea, send a DELETE request to `/tea/{teaID}/owner/{ownerID}`
