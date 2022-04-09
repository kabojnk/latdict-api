# Latdict API

## Intro and Tech Stack
- This API is written in [Go](https://go.dev/), using the [Gin framework](https://github.com/gin-gonic/gin)
- The main data source is [PostgreSQL](https://www.postgresql.org/)
- Caching utilizes [redis](https://redis.io/)
- Web server is `nginx`, with SSL via `letsencrypt`
- A `docker-compose.dev.yml` file is supplied to spin up local postgres, redis, nginx and nginx-letsencrypt containers for local development

## Building

Site uses environmental variables, which can optionally be supplied in a `.env` file for local development.
There is a `.env.example` file that can be copied to a new `.env` file and filled out with the various values you want 
(like database credentials).

1. Ideally, you would have **Go 1.16+** installed.
2. Run `go run cmd/.` (`main.go` is inside the `cmd/` directory)
3. If you wish to do hot reloading, the API supports [air](https://github.com/cosmtrek/air). Once installed, just run `air` in the project root.

## Swagger API documentation
 
The API uses [gin-swagger](https://github.com/swaggo/gin-swagger).

### Viewing documentation

1. Run the server
2. Go to http://localhost:8000/swagger/index.html in a browser to view the API documentation

### Generating (or regenerating) documentation

Documentation is based primarily on code comments that define the various endpoint handlers. To generate new documentation:

1. Install [Swag for Go](https://github.com/swaggo/swag): `go get -u github.com/swaggo/swag/cmd/swag`
2. Type `swag init --parseDependency` -- the `--parseDependency` flag is important so that it can find the typedef for the response models

## Notes on types

- Models in Go are typically `struct` typedefs, and their fields are inherently private if they begin with a lowercase letter. So we start each field with an uppercase letter, and just explicitly specify serialization mapping in the typedef itself. It's a bit verbose but it's not the end of the world. 
- In `types/types.go` directory you might see some models that look redundant. E.g. `DBEntry` and `APIEntry`. Any DB model 
that has a `Sql.Null*` type in its field list (like `Sql.NullString`), will have a corresponding API Model that will have
those values unboxed, so you don't get weird values to unbox like `Entry.AdditionalInfo.Age.String`. Plus having the 
extra context bounds makes it easier to dictate what ultimately gets pushed to the API output.


