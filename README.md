# URL shortening

<!--toc:start-->
- [URL shortening](#url-shortening)
  - [API Endpoints:](#api-endpoints)
<!--toc:end-->

This URL shortening service, built with Go and vertical slice  architecture for efficient scalability and performance.

## API Endpoints

- Shorten URL

```
POST /api/shorten
```

Request: { "originalUrl": "<https://example.com/long-path>", "customAlias": "mylink" (optional) }
Response: { "shortUrl": "<https://short.io/abc123>", "originalUrl": "<https://example.com/long-path>", "createdAt": "2025-05-05T10:30:00Z", "expiresAt": "2026-05-05T10:30:00Z" (optional) }

- Redirect

```
GET /{shortCode}
```

Action: 302 redirect to original URL + record visit

- Get URL Info

```
GET /api/urls/{shortCode}
```

Response: { "shortUrl": "<https://short.io/abc123>", "originalUrl": "<https://example.com/long-path>", "createdAt": "2025-05-05T10:30:00Z", "visits": 42 }

- Delete URL

```
DELETE /api/urls/{shortCode}
```

Response: { "success": true }

Step 1:

make download-migrate
