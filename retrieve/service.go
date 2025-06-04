package retrieve

import (
	"context"
	"log"
)

type Store interface {
	RertrieveShortUrl(ctx context.Context, shortUrl string) (string, error)
}
type Cache interface {
	Save(ctx context.Context, shortUrl, originalUrl string) error
	Get(ctx context.Context, shortUrl string) (string, error)
}

type RetrieveService struct {
	store Store
	cache Cache
}

func NewRetrieveService(store Store, cache Cache) *RetrieveService {
	return &RetrieveService{
		store: store,
		cache: cache,
	}
}

func (r *RetrieveService) RetrieveShortUrl(ctx context.Context, shortUrl string) (string, error) {
	originalUrl, err := r.cache.Get(ctx, shortUrl)
	if err != nil {
		return "", err
	}
	if originalUrl == "" {

		originalUrl, err = r.store.RertrieveShortUrl(ctx, shortUrl)
		if err != nil {
			return "", err
		}

		log.Printf("shortUrl %s got from the database\n", shortUrl)

	} else {
		log.Printf("shortUrl %s got from the cache\n", shortUrl)
	}
	err = r.cache.Save(ctx, shortUrl, originalUrl)
	if err != nil {
		return "", err
	}
	return originalUrl, nil
}
