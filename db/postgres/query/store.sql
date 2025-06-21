-- name: RetrieveShortUrl :one
SELECT original_url, expired_at FROM urls
WHERE short_url = $1;

-- name: InsertNewShortUrl :one
INSERT INTO urls (original_url, short_url, expired_at) VALUES ( $1, $2, $3) RETURNING short_url;

-- name: CheckShortUrlExists :one
SELECT EXISTS (
    SELECT 1 FROM urls
    WHERE short_url = $1
);

-- name: GetInfoUrl :one
SELECT original_url FROM urls
WHERE short_url = $1;

-- name: DeleteShortUrl :exec
DELETE FROM urls WHERE short_url = $1;

