package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	db "url-shortening-service/db/postgres/sqlc"
	"url-shortening-service/internal/repository"
)

type PostgresStore struct {
	db      *sql.DB
	queries *db.Queries
}

// DeleteShortUrl implements repository.Store.
func (g *PostgresStore) DeleteShortUrl(ctx context.Context, shortUrl string) error {
	err := g.queries.DeleteShortUrl(ctx, shortUrl)
	if err != nil {
		return fmt.Errorf("cannot delete shortUrl from database %w", err)
	}
	return nil
}

// GetInfoUrl implements repository.Store.
func (g *PostgresStore) GetInfoUrl(ctx context.Context, shortUrl string) (string, error) {
	panic("unimplemented")
}

func NewPostgresStore(dbConn *sql.DB) *PostgresStore {
	return &PostgresStore{
		db:      dbConn,
		queries: db.New(dbConn),
	}
}

func (g *PostgresStore) InsertNewShortUrl(ctx context.Context, originalUrl, shortUrl string) (string, time.Time, error) {
	expiredAt := time.Now().Add(5 * time.Hour)
	params := db.InsertNewShortUrlParams{
		OriginalUrl: originalUrl,
		ShortUrl:    shortUrl,
		ExpiredAt:   expiredAt,
	}

	url, err := g.queries.InsertNewShortUrl(ctx, params)
	if err != nil {
		return "", time.Time{}, err
	}
	return url, expiredAt, err
}

func (p *PostgresStore) RertrieveShortUrl(ctx context.Context, shortUrl string) (*repository.RetrieveShortUrlRow, error) {
	url, err := p.queries.RetrieveShortUrl(ctx, shortUrl)
	if err != nil {
		return &repository.RetrieveShortUrlRow{}, repository.ErrNotFound
	}
	return &repository.RetrieveShortUrlRow{
		OriginalUrl: url.OriginalUrl,
		ExpiredAt:   url.ExpiredAt,
	}, nil
}

func (p *PostgresStore) CheckShortUrlExists(ctx context.Context, shortUrl string) (bool, error) {
	shortUrlExist, err := p.queries.CheckShortUrlExists(ctx, shortUrl)
	if err != nil {
		return false, err
	}
	return shortUrlExist, nil
}
