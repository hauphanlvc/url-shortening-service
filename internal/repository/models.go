package repository

import (
	"context"
	"errors"
	"time"
)

type RetrieveShortUrlRow struct {
	OriginalUrl string
	ExpiredAt   time.Time
}

var ErrNotFound = errors.New("url not found")

type Store interface {
	InsertNewShortUrl(ctx context.Context, originalUrl, shortUrl string) (string, time.Time, error)
	RertrieveShortUrl(ctx context.Context, shortUrl string) (*RetrieveShortUrlRow, error)
	CheckShortUrlExists(ctx context.Context, shortUrl string) (bool, error)
	DeleteShortUrl(ctx context.Context, shortUrl string) error
	GetInfoUrl(ctx context.Context, shortUrl string) (string, error)
}
