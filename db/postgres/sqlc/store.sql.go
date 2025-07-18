// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: store.sql

package db

import (
	"context"
	"time"
)

const checkShortUrlExists = `-- name: CheckShortUrlExists :one
SELECT EXISTS (
    SELECT 1 FROM urls
    WHERE short_url = $1
)
`

func (q *Queries) CheckShortUrlExists(ctx context.Context, shortUrl string) (bool, error) {
	row := q.db.QueryRowContext(ctx, checkShortUrlExists, shortUrl)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const deleteShortUrl = `-- name: DeleteShortUrl :exec
DELETE FROM urls WHERE short_url = $1
`

func (q *Queries) DeleteShortUrl(ctx context.Context, shortUrl string) error {
	_, err := q.db.ExecContext(ctx, deleteShortUrl, shortUrl)
	return err
}

const getInfoUrl = `-- name: GetInfoUrl :one
SELECT original_url FROM urls
WHERE short_url = $1
`

func (q *Queries) GetInfoUrl(ctx context.Context, shortUrl string) (string, error) {
	row := q.db.QueryRowContext(ctx, getInfoUrl, shortUrl)
	var original_url string
	err := row.Scan(&original_url)
	return original_url, err
}

const insertNewShortUrl = `-- name: InsertNewShortUrl :one
INSERT INTO urls (original_url, short_url, expired_at) VALUES ( $1, $2, $3) RETURNING short_url
`

type InsertNewShortUrlParams struct {
	OriginalUrl string
	ShortUrl    string
	ExpiredAt   time.Time
}

func (q *Queries) InsertNewShortUrl(ctx context.Context, arg InsertNewShortUrlParams) (string, error) {
	row := q.db.QueryRowContext(ctx, insertNewShortUrl, arg.OriginalUrl, arg.ShortUrl, arg.ExpiredAt)
	var short_url string
	err := row.Scan(&short_url)
	return short_url, err
}

const retrieveShortUrl = `-- name: RetrieveShortUrl :one
SELECT original_url, expired_at FROM urls
WHERE short_url = $1
`

type RetrieveShortUrlRow struct {
	OriginalUrl string
	ExpiredAt   time.Time
}

func (q *Queries) RetrieveShortUrl(ctx context.Context, shortUrl string) (RetrieveShortUrlRow, error) {
	row := q.db.QueryRowContext(ctx, retrieveShortUrl, shortUrl)
	var i RetrieveShortUrlRow
	err := row.Scan(&i.OriginalUrl, &i.ExpiredAt)
	return i, err
}
