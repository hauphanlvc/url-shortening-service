package retrieve

import (
	"context"
	_ "log"
	"url-shortening-service/internal/cache"
	"url-shortening-service/internal/repository"

	"github.com/rs/zerolog/log"
)

//	type Store interface {
//		RertrieveShortUrl(ctx context.Context, shortUrl string) (string, error)
//	}
//
//	type Cache interface {
//		Save(ctx context.Context, shortUrl, originalUrl string) error
//		Get(ctx context.Context, shortUrl string) (string, error)
//	}
type RetrieveService struct {
	store repository.Store
	cache cache.Cache
}

func NewRetrieveService(store repository.Store, cache cache.Cache) *RetrieveService {
	return &RetrieveService{
		store: store,
		cache: cache,
	}
}

func (r *RetrieveService) RetrieveShortUrl(ctx context.Context, shortUrl string) (string, error) {
	log.Info().Msgf("RetrieveShortUrl service")
	result, err := r.cache.Get(ctx, shortUrl)
	if err != nil {
		log.Error().Err(err)
		return "", err
	}
	if result.OriginalUrl != "" {

		log.Debug().Msgf("Cache HIT")
		log.Debug().Str("shoshortUrl", shortUrl).Str("originalUrl", result.OriginalUrl).Msgf("shortUrl %s got from the cache\n", shortUrl)
	} else {

		log.Debug().Msgf("Cache miss")
		// originalUrl, err = r.store.RertrieveShortUrl(ctx, shortUrl)
		result, err = r.store.RertrieveShortUrl(ctx, shortUrl)
		if err != nil {

			log.Error().Err(err)
			return "", err
		}

		log.Debug().Str("shortUrl", shortUrl).Msgf("shortUrl %s got from the database\n", shortUrl)

		err = r.cache.Save(ctx, shortUrl, repository.RetrieveShortUrlRow{OriginalUrl: result.OriginalUrl, ExpiredAt: result.ExpiredAt})
		if err != nil {
			log.Error().Err(err)
			return "", err
		}
	}
	return result.OriginalUrl, nil
}
