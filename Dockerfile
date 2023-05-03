FROM golang:1.20-alpine as builder

WORKDIR /go/src/app
COPY . .
RUN go get ./...
RUN go build -o /go/bin/ ./...

FROM alpine:latest
COPY --from=builder /go/bin/urlshortener /app
ENV GIN_MODE=release
ENTRYPOINT /app

EXPOSE 8000
