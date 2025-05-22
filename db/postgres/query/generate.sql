-- name: RetrieveShortUrl :one
SELECT * FROM urls
WHERE short_url = $1 LIMIT 1;

-- name: InsertNewShortUrl :one
INSERT INTO urls (original_url, short_url, expired_at) VALUES ( $1, $2, $3) RETURNING *;

-- name: DeleteExpiredShortUrl :exec
DELETE FROM urls WHERE short_url = $1;

-- name: CheckShortUrlExists :one
SELECT EXISTS (
    SELECT 1 FROM urls
    WHERE short_url = $1
);
