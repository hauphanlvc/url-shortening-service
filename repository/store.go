package repository

import (
	"context"
	"database/sql"
	"time"
	db "url-shortening-service/db/postgres/sqlc"
)

type Store interface {
	InsertNewShortUrl(ctx context.Context, originalUrl, shortUrl string) (*db.Url, error)
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

func (g *PostgresStore) InsertNewShortUrl(ctx context.Context, originalUrl, shortUrl string) (*db.Url, error) {
	params := db.InsertNewShortUrlParams{
		OriginalUrl: originalUrl,
		ShortUrl:    shortUrl,
		ExpiredAt:   time.Now().Add(5 * time.Hour),
	}

	url, err := g.queries.InsertNewShortUrl(ctx, params)
	if err != nil {
		return nil, err
	}
	return &url, err
}
