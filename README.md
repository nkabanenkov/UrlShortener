# An URL shortener service

- [x] Docker image, compose & [devcontainer](https://containers.dev) files
- [x] Endpoints:
  - [x] HTTP
  - [x] gRPC
- [x] In memory & PostgreSQL storages
- [x] Tests:
    - [x] base encoder
    - [x] URL characters (RFC3986-approved characters) and scheme (http or https) validators
    - [x] in-memory storage
    - [x] PostgreSQL storage (with [testcontainers](https://golang.testcontainers.org))

## Usage

### Running the service

- Edit `.env` if needed,
- `cd` to the root of the repository, then do `docker compose up`.

If you want to run in-memory storage, then do `docker compose -f docker-compose.inmemory.yml up`.

**The service's HTTP endpoints will be available at `8000`, gRPC endpoints will be available at `9111`.**

### HTTP endpoints
- POST `/urls`, where body of the request is in form `{"url": "https://google.com"}`, transforms the given URL to a shortened form and returns it.

Example request:
```
curl -X "POST" "http://127.0.0.1:8000/urls" \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{
  "url": "https://google.com"
}'
```
Example answer:
```
HTTP/1.1 200 OK
...
{"message":"eaaaaaaaaa"}
```

- GET `/urls/[url]`, where `[url]` is the shortened url, converts the `[url]` to full address.

Example request:
```
curl "http://127.0.0.1:8000/urls/eaaaaaaaaa"
```

Example answer:
```
HTTP/1.1 200 OK
...
{"message":"https://google.com"}
```

### gRPC endpoints
See `api/urlshortener.proto`.

### Encoding overflow
When the encoder can't fit the id of an URL into a word encoded with the given alphabet and length, `EncodingOverflowError{"encoding overflow"}` is returned and the service exits with non-zero code.

### Testing
PostgreSQL storage tests are performed with [testcontainers](https://golang.testcontainers.org).

### Developing
VS Code [devcontainer](https://containers.dev) files are located in `.devcontainer`. Press Cmd+Shift+P and select `Reopen in container`.
