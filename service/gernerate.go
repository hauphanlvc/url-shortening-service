package service

import (
	"context"
	"database/sql"
	"fmt"
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
	shortUrl, err := GenerateShortUrl(url)
	if err != nil {
		return nil, err
	}
	salt := 0
	for true {

		existShortUrl, err := gs.queries.CheckShortUrlExists(ctx, shortUrl)
		if err != nil {
			return nil, err
		}
		if existShortUrl {
			shortUrl, err = GenerateShortUrl(fmt.Sprintf("%s%d", shortUrl, salt))
			if err != nil {
				return nil, err
			}
			salt += 1
		} else {
			break
		}
	}
	params := db.InsertNewShortUrlParams{
		OriginalUrl: url,
		ShortUrl:    shortUrl,
		ExpiredAt:   time.Now().Add(5 * time.Hour),
	}
	result, err := gs.queries.InsertNewShortUrl(ctx, params)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
