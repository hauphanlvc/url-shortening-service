# URL shortening

<!--toc:start-->
- [URL shortening](#url-shortening)
  - [API Endpoints](#api-endpoints)
  - [Diagram architecture](#diagram-architecture)
  - [To do](#to-do)
<!--toc:end-->

- [API Endpoints](#api-endpoints)
- [System design architecture](#diagram-architecture)
- [Todo](#to-do)
<!--toc:end-->

This URL shortening service, built with Go and vertical slice  architecture for
efficient scalability and performance.

## API Endpoints

- Shorten URL

```
POST /shorten
```

Request: {
    "originalUrl": "<https://example.com/long-path>",
}

Response: {
    "shortUrl": "<https://short.io/abc123>",
}

- Redirect

```
GET /{shortCode}
```

Action: 302 redirect to original URL

## Diagram architecture

## To do

- [x] Add logging
- [x] Delete expired short link
- [x] Use dragonfly for caching
