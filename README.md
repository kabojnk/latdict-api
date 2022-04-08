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
3. If you wish to do hot reloading, the API supports `[air](https://github.com/cosmtrek/air)`. Once installed, just run `air` in the project root.

## Notes on types

- Models in Go are typically `struct` typedefs, and their fields are inherently private if they begin with a lowercase letter. So we start each field with an uppercase letter, and just explicitly specify serialization mapping in the typedef itself. It's a bit verbose but it's not the end of the world. 
- In `types/types.go` directory you might see some models that look redundant. E.g. `DBEntry` and `APIEntry`. Any DB model 
that has a `Sql.Null*` type in its field list (like `Sql.NullString`), will have a corresponding API Model that will have
those values unboxed, so you don't get weird values to unbox like `Entry.AdditionalInfo.Age.String`. Plus having the 
extra context bounds makes it easier to dictate what ultimately gets pushed to the API output.


