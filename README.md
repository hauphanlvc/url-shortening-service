# ✂️ URL Shortening Service

A high-performance, minimalistic URL shortening service built with Go, PostgreSQL, Redis (DragonFly), and Gin.

## ✨ Features

- 🔗 **URL Shortening**: Generate short, unique URLs using NanoID.
- 🚀 **Fast Redirection**: Retrieve and redirect using in-memory Redis cache (DragonFly).
- 💾 **PostgreSQL Persistence**: Durable storage of original-short URL mappings.
- 🧠 **Smart Caching**: Redis-first lookups for optimal performance.
- ⏳ **Rate Limiting**: Basic rate limiter to prevent abuse (to be improved per client IP).
- 🧱 **Vertical Architecture**: Modular codebase using services, transport layers, and repositories.

---

## 📁 Project Structure

```
.
├── cmd/apis/         # Main application entrypoint
├── internal/         # Internal application logic
│   ├── cache/        # Redis/Dragonfly caching logic
│   ├── repository/   # Database interactions (PostgreSQL)
│   └── urls/         # Core URL shortening feature domain
│       ├── generate/ # Logic for generating short URLs
│       └── retrieve/ # Logic for retrieving original URLs
├── config/           # Configuration loading
├── db/postgres/      # Database migration files for Postgresql
├── go.mod
└── Makefile
└── sqlc.yaml         # generated sqlc code

```

---

## ⚙️ Requirements

- Go 1.21+
- Docker compose
- PostgreSQL
- Redis (preferably [DragonFly](https://www.dragonflydb.io/))
- [golang-migrate/migrate](https://github.com/golang-migrate/migrate)
- [air-verse/air](https://github.com/air-verse/air)

---

## 🚀 Getting Started

### 1. Clone and Set Up Environment
```bash
git clone [https://github.com/yourname/url-shortening-service.git](https://github.com/yourname/url-shortening-service.git)
cd url-shortening-service
cp .env.example .env
```
*After cloning, update the `.env` file with your PostgreSQL and Redis connection details.*

### 2. Set Up the Database
*Ensure Docker is running, then start the database containers and apply migrations.*
```bash
# Start required database containers (e.g., via docker-compose)
make db-up

# Apply database migrations
make migrate-up
```

### 3. Run the Application
*You can run the application using one of the following commands:*
```bash
# To run with live-reloading (recommended for development)
air

# Or to run using the makefile
make run

# Or to run directly with Go
go run cmd/apis/main.go
```

---

## 🌐 API Endpoints

The API is built using Gin and provides endpoints for creating and retrieving short URLs.

A global rate limiter is applied to all endpoints, currently set to **1 request per second** with a burst capacity of **4**. If the limit is exceeded, the API will respond with `429 Too Many Requests`.

---

### Create a Short URL

Creates a new short URL and stores it.

* **Endpoint:** `POST /shorten`
* **Request Body:**

    ```json
    {
      "url": "https://github.com/hauphanlvc/url-shortening-service"
    }
    ```

* **Success Response (`200 Created`):**

    The response body includes the full short URL.

    ```json
    {
      "ShortUrl": "o5CiXqw"
    }
    ```

* **Error Responses:**
    * `400 Bad Request`: If the `long_url` is missing or invalid. (Implementation Future)
    * `429 Too Many Requests`: If the rate limit is exceeded.
    * `500 Internal Server Error`: For any unexpected server-side errors.

* **Example `curl`:**

    ```bash
    curl -X POST http://localhost:8080/shorten \
    -H "Content-Type: application/json" \
    -d '{"url": "https://github.com/hauphanlvc/url-shortening-service"}'
    ```

---

### Redirect to Original URL

Retrieves the original long URL from a short code and performs an HTTP redirect. The service first checks the cache; on a cache miss, it queries the database and populates the cache.

* **Endpoint:** `GET /:shortUrl`
* **Description:** This endpoint is the primary redirection mechanism. Accessing this URL in a browser will redirect the user to the original long URL.
* **Success Response (`302 Found`):**
    * An HTTP redirect to the original `long_url`. The `Location` header will contain the destination URL.

* **Error Responses:**
    * `404 Not Found`: If the `shortUrl` code does not exist in the database.
    * `429 Too Many Requests`: If the rate limit is exceeded.

* **Example `curl`:**

    Use the `-L` flag to make curl follow the redirect.

    ```bash
    curl -L http://localhost:8080/jA8s1bC
    ```

---

### Get URL Details (Future Implementation)

This endpoint is defined in the router but does not have an implementation yet. It's intended to fetch details about a short URL without performing a redirect.

* **Endpoint:** `GET /urls/:shortUrl`
* **Success Response (`200 OK`):**

    *(Anticipated Response)*
    ```json
    {
      "long_url": "[https://example.com/a-very-long-and-unwieldy-url-to-shorten](https://example.com/a-very-long-and-unwieldy-url-to-shorten)",
      "short_url": "http://localhost:8080/jA8s1bC",
      "created_at": "2025-07-13T14:57:31Z"
    }
    ```

---

### Delete a URL (Future Implementation)

This endpoint is defined in the router but does not have an implementation yet. It's intended to delete a short URL record from the system.

* **Endpoint:** `DELETE /urls/:shortUrl`
* **Success Response (`204 No Content`):**
    * An empty response body indicating successful deletion.

* **Error Responses:**
    * `404 Not Found`: If the `shortUrl` code does not exist.
