package retrieve

import (
	"context"
	_ "log"
	"time"
	"url-shortening-service/internal/cache"
	"url-shortening-service/internal/repository"

	"github.com/rs/zerolog/log"
)

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
	if result != nil {

		log.Debug().Msgf("Cache HIT")
		log.Debug().Str("shortUrl", shortUrl).Str("originalUrl", result.OriginalUrl).Msgf("shortUrl %s got from the cache\n", shortUrl)
		if result.ExpiredAt.Before(time.Now()) {
			go func() {

				err := r.cache.Delete(ctx, shortUrl)
				if err != nil {
					log.Error().Err(err).Str("shortUrl", shortUrl).Msg("failed to delete expired shortUrl asynchronously from cache")
				}

				log.Debug().Str("shortUrl", shortUrl).Msgf("delete shortUrl succesful from cache")
			}()
			return "", repository.ErrNotFound
		}

	} else {

		log.Debug().Msgf("Cache miss")
		result, err = r.store.RertrieveShortUrl(ctx, shortUrl)
		if err != nil {
			log.Error().Err(err)
			return "", err
		}

		if result.ExpiredAt.Before(time.Now()) {
			go func() {

				log.Debug().Msgf("the shortUrl have expired, we gonna to delete it")
				err := r.store.DeleteShortUrl(context.Background(), shortUrl)

				if err != nil {
					log.Error().Err(err).Str("shortUrl", shortUrl).Msg("failed to delete expired shortUrl asynchronously from database\n")
				}
				log.Debug().Str("shortUrl", shortUrl).Msgf("delete shortUrl succesful from database\n")

			}()

			return "", repository.ErrNotFound
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
