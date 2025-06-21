package generate

import (
	"context"
	"errors"
	_ "log"
	"url-shortening-service/internal/cache"
	"url-shortening-service/internal/repository"

	"github.com/rs/zerolog/log"
)

const MAX_RETRIES = 1000000

// type Store interface {
// 	InsertNewShortUrl(ctx context.Context, originalUrl, shortUrl string) (string, error)
// 	RertrieveShortUrl(ctx context.Context, shortUrl string) (, error)
// 	CheckShortUrlExists(ctx context.Context, shortUrl string) (bool, error)
// 	DeleteShortUrl(ctx context.Context, shortUrl string) error
// 	GetInfoUrl(ctx context.Context, shortUrl string) (string, error)
// }

var ErrCannotGernerateShortURL = errors.New("cannot generate short url, always got the collision")

// type Cache interface {
// 	Save(ctx context.Context, shortUrl, originalUrl string) error
// 	Get(ctx context.Context, shortUrl string) (string, error)
// }

type GenerateService struct {
	store     repository.Store
	generator Generator
	cache     cache.Cache
}

func NewGenerateService(store repository.Store, generator Generator, cache cache.Cache) *GenerateService {
	return &GenerateService{
		store:     store,
		generator: generator,
		cache:     cache,
	}
}

func (g *GenerateService) ShortUrlExist(ctx context.Context, shortUrl string) (bool, error) {
	log.Info().Msg("ShortUrlExist function")
	shortUrlExist, err := g.store.CheckShortUrlExists(ctx, shortUrl)
	if err != nil {
		log.Error().Err(err)
		return false, err
	}
	log.Debug().
		Str("shortUrl", shortUrl).
		Bool("shoshortUrlExist", shortUrlExist)
	return shortUrlExist, nil
}

func (g *GenerateService) GetUniqueShortUrl(ctx context.Context) (string, error) {
	log.Info().Msg("GetUniqueShortUrl function")
	counter := 1
	var shortUrl string
	var err error
	for {

		shortUrl, err = g.generator.GenerateShortUrl()
		if err != nil {
			return "", err
		}

		exists, err := g.ShortUrlExist(ctx, shortUrl)
		if err != nil {
			return "", err
		}
		if !exists {
			break
		}

		log.Debug().Msgf("The short url exists, retry %d\n", counter)
		if counter > MAX_RETRIES {
			log.Error().Err(ErrCannotGernerateShortURL)
			return "", ErrCannotGernerateShortURL
		}

		counter += 1
		continue
	}
	log.Debug().Str("shortUrl", shortUrl)
	return shortUrl, nil
}

func (g *GenerateService) InsertNewShortUrl(ctx context.Context, originalUrl string) (string, error) {
	log.Info().Msg("InsertNewShortUrl")
	shortUrl, err := g.GetUniqueShortUrl(ctx)
	log.Debug().Str("shortUrl", shortUrl)
	if err != nil {
		log.Error().Err(err)
		return "", err
	}
	shortUrlResult, expiredAt, err := g.store.InsertNewShortUrl(ctx, originalUrl, shortUrl)
	if err != nil {
		log.Error().Err(err)
		return "", err
	}
	log.Debug().Str("result", shortUrlResult).Msg("result after inserting new short url to database")
	// err = g.cache.Save(ctx, shortUrl, originalUrl)
	err = g.cache.Save(ctx, shortUrl, repository.RetrieveShortUrlRow{OriginalUrl: originalUrl, ExpiredAt: expiredAt})
	if err != nil {
		log.Error().Err(err)
		return "", err
	}
	return shortUrlResult, nil
}
