package retrieve

import (
	"context"
	_ "log"

	"github.com/rs/zerolog/log"
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
	log.Info().Msgf("RetrieveShortUrl service")
	originalUrl, err := r.cache.Get(ctx, shortUrl)
	if err != nil {
		log.Error().Err(err)
		return "", err
	}
	if originalUrl == "" {

		log.Debug().Msgf("Cache miss")
		originalUrl, err = r.store.RertrieveShortUrl(ctx, shortUrl)
		if err != nil {

			log.Error().Err(err)
			return "", err
		}

		log.Debug().Str("shortUrl", shortUrl).Msgf("shortUrl %s got from the database\n", shortUrl)

	} else {

		log.Debug().Msgf("Cache HIT")
		log.Debug().Str("shoshortUrl", shortUrl).Str("originalUrl", originalUrl).Msgf("shortUrl %s got from the cache\n", shortUrl)
	}
	err = r.cache.Save(ctx, shortUrl, originalUrl)
	if err != nil {
		log.Error().Err(err)
		return "", err
	}
	return originalUrl, nil
}
