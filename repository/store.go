package repository

import (
	"context"
	"database/sql"
	"time"
	db "url-shortening-service/db/postgres/sqlc"
)

type Store interface {
	InsertNewShortUrl(ctx context.Context, originalUrl, shortUrl string) (string, error)
	RertrieveShortUrl(ctx context.Context, shortUrl string) (*string, error)
	CheckShortUrlExists(ctx context.Context, shortUrl string) (bool, error)
}

type PostgresStore struct {
	db      *sql.DB
	queries *db.Queries
}

func NewPostgresStore(dbConn *sql.DB) *PostgresStore {
	return &PostgresStore{
		db:      dbConn,
		queries: db.New(dbConn),
	}
}

func (g *PostgresStore) InsertNewShortUrl(ctx context.Context, originalUrl, shortUrl string) (string, error) {
	params := db.InsertNewShortUrlParams{
		OriginalUrl: originalUrl,
		ShortUrl:    shortUrl,
		ExpiredAt:   time.Now().Add(5 * time.Hour),
	}

	url, err := g.queries.InsertNewShortUrl(ctx, params)
	if err != nil {
		return "", err
	}
	return url.ShortUrl, err
}

func (p *PostgresStore) RertrieveShortUrl(ctx context.Context, shortUrl string) (string, error) {
	url, err := p.queries.RetrieveShortUrl(ctx, shortUrl)
	if err != nil {
		return "", err
	}
	return url.OriginalUrl, nil
}

func (p *PostgresStore) CheckShortUrlExists(ctx context.Context, shortUrl string) (bool, error) {
	shortUrlExist, err := p.queries.CheckShortUrlExists(ctx, shortUrl)
	if err != nil {
		return false, err
	}
	return shortUrlExist, nil
}
