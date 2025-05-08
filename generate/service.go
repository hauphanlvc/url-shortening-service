package generate

import (
	"context"
	"database/sql"
	"time"
	db "url-shortening-service/db/sqlc"
)

type GenerateService struct {
	db      *sql.DB
	queries *db.Queries
}

func NewGenerateService(dbConn *sql.DB) *GenerateService {
	return &GenerateService{
		db:      dbConn,
		queries: db.New(dbConn),
	}
}

func (gs *GenerateService) InsertNewShortUrl(ctx context.Context, url string) (*db.Url, error) {
	shortUrl, err := generateShortUrl(url)
	if err != nil {
		return nil, err
	}
	params := db.InsertNewShortUrlParams{
		OriginalUrl: url,
		ShortUrl:    shortUrl,
		ExpiredAt:   time.Now(),
	}
	result, err := gs.queries.InsertNewShortUrl(ctx, params)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
