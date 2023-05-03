# An URL shortener service

- [x] Docker image, compose & devcontainer files
- [x] In memory & PostgreSQL storages
- [x] Tests:
    - [x] base encoder
    - [x] http prefix validator
    - [x] in-memory storage
    - [x] PostgreSQL (with [testcontainers](https://golang.testcontainers.org))

## Usage

### Running the service

- Edit `.env` if needed,
- `cd` to the root of the repository, then do `docker compose up`. The service will be listening on `8000`.

If you want to run in-memory storage, then do `docker compose -f docker-compose.inmemory.yml up`.

### REST methods
- GET `/urls/[url]`, where `[url]` is the shortened url, converts the `[url]` to full address,
- POST `/urls`, where body of the request is in form `{"url": "https://google.com"}`, transforms the given URL to a shortened form and returns it.

### Encoding overflow
When the encoder can't fit the id of an URL into a word encoded with the given alphabet and length, `EncodingError{"encoding overflow"}` is returned and the service panics.

### Testing
PostgreSQL storage tests are performed with [testcontainers](https://golang.testcontainers.org).

### Developing
VS Code [devcontainer](https://containers.dev) files are located in `.devcontainer`. Press Cmd+Shift+P and select `Reopen in container`.
