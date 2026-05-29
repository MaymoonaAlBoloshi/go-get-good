# express-go

Tiny HTTP/1.1 server in Go.

This was written for learning purposes, mostly through the CodeCrafters HTTP
server challenge. It is not a production server. Some parts are rough, direct,
and incomplete, because the point was to learn how HTTP works by building the
pieces by hand.

## What exists in it

- a TCP server listening on port `4221`
- a tiny router with `GET` and `POST`
- basic HTTP request parsing
- HTTP responses with status, headers, and body
- routes for `/`, `/echo/*`, `/user-agent`, and `/files/*`
- file reads and writes with `--directory`
- gzip response compression for `/echo/*`
- persistent connections unless `Connection: close` is sent

## Run it

```bash
go run ./app
```

With a file directory:

```bash
go run ./app --directory /tmp/
```

For development with `air`:

```bash
./dev.sh
```

## Routes

`GET /`
Returns `200 OK`.

```bash
curl -i http://localhost:4221/
```

`GET /echo/*`
Echoes the path text back as `text/plain`.

```bash
curl -i http://localhost:4221/echo/hello
```

`GET /user-agent`
Returns the `User-Agent` header.

```bash
curl -i http://localhost:4221/user-agent
```

`GET /files/*`
Reads a file from the directory passed with `--directory`.

```bash
curl -i http://localhost:4221/files/test.txt
```

`POST /files/*`
Writes the request body to a file in the directory passed with `--directory`.

```bash
curl -i -X POST http://localhost:4221/files/test.txt -d "hello"
```

## Compression

Supported:

```bash
curl -i -H "Accept-Encoding: gzip" http://localhost:4221/echo/hello
```

## Limits

- not a real web framework
- no TLS
- no middleware
- no query parsing
- no chunked requests or responses
- request parsing is simple and assumes small requests
- wildcard routing only supports routes ending in `/*`
- file paths are intentionally simple and not hardened
