package retrieve

import (
	"context"
	"url-shortening-service/repository"
)

type RetrieveService struct {
	store repository.Store
}

func NewRetrieveService(store repository.Store) *RetrieveService {
	return &RetrieveService{
		store: store,
	}
}

func (r *RetrieveService) RetrieveShortUrl(ctx context.Context, shortUrl string) (*string, error) {
	originalUrl, err := r.store.RertrieveShortUrl(ctx, shortUrl)
	if err != nil {
		return nil, err
	}
	return originalUrl, nil
}
