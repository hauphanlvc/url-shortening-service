# URL shortening

- [API Endpoints](#api-endpoints)
- [System design architecture](#diagram-architecture)
- [Todo](#to-do)
<!--toc:end-->

This URL shortening service, built with Go and vertical slice  architecture for
efficient scalability and performance.

## API Endpoints

- Shorten URL

```bash
POST /shorten
```

Request: {
    "originalUrl": "<https://example.com/long-path>",
    "customAlias": "mylink" (optional) // To do
}

Response: {
    "shortUrl": "<https://short.io/abc123>",
    <!-- "originalUrl": "<https://example.com/long-path>", -->
    <!-- "createdAt": "2025-05-05T10:30:00Z", -->
    <!-- "expiresAt": "2026-05-05T10:30:00Z" (optional) -->
}

- Redirect

```bash
GET /{shortCode}
```

Action: 302 redirect to original URL

- Get URL Info

```bash
GET /urls/{shortCode}
```

Response: {
    "shortUrl": "<https://short.io/abc123>",
    "originalUrl": "<https://example.com/long-path>",
    "createdAt": "2025-05-05T10:30:00Z",
    "expiresAt": "2026-05-05T10:30:00Z"
    "visits": 42
}

- Delete URL

```bash
DELETE /urls/{shortCode}
```

Response: { "success": true }

## Diagram architecture

## To do
- [ ] Add logging
- [ ] Delete expired short link
- [ ] Use valkey for caching
- [ ] Add user id into URL table and authenticate user by Oauth2
- [ ] Custom alias feature ( limit: 7 characters or more)
