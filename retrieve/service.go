package retrieve

import (
	"context"
	"database/sql"
	db "url-shortening-service/db/postgres/sqlc"
)

type RetrieveService struct {
	db      *sql.DB
	queries *db.Queries
}

func NewRetrieveService(dbConn *sql.DB) *RetrieveService {
	return &RetrieveService{
		db:      dbConn,
		queries: db.New(dbConn),
	}
}

func (rs *RetrieveService) RetrieveShortUrl(ctx context.Context, shortUrl string) (*db.Url, error) {
	OriginalUrl, err := rs.queries.RetrieveShortUrl(ctx, shortUrl)
	if err != nil {
		return nil, err
	}
	return &OriginalUrl, nil
}
