# Proxy Server

A simple caching reverse proxy server built with Go.

## How it works

Requests are forwarded to the target URL and cached for one hour. Subsequent requests to the same URL are served from cache without hitting the upstream server.

## Running with Docker
```bash
docker build -t proxy-server .
docker run -p 8080:8080 proxy-server
```

## Running locally
```bash
go run .
```

## Usage
```bash
curl "http://localhost:8080/proxy?url=https://example.com"
```

Check the `X-Cache` header in the response:
```
X-Cache: MISS  → fetched from upstream
X-Cache: HIT   → served from cache
```

## Cache

- TTL: 1 hour
- Expired entries are cleaned up every 10 minutes