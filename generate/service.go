package generate

import (
	"context"
	"errors"
	_ "log"

	"github.com/rs/zerolog/log"
)

const MAX_RETRIES = 1000000

type Store interface {
	InsertNewShortUrl(ctx context.Context, originalUrl, shortUrl string) (string, error)
	RertrieveShortUrl(ctx context.Context, shortUrl string) (string, error)
	CheckShortUrlExists(ctx context.Context, shortUrl string) (bool, error)
}

var ErrCannotGernerateShortURL = errors.New("cannot generate short url, always got the collision")

type Cache interface {
	Save(ctx context.Context, shortUrl, originalUrl string) error
	Get(ctx context.Context, shortUrl string) (string, error)
}

type GenerateService struct {
	store     Store
	generator Generator
	cache     Cache
}

func NewGenerateService(store Store, generator Generator, cache Cache) *GenerateService {
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
	result, err := g.store.InsertNewShortUrl(ctx, originalUrl, shortUrl)
	if err != nil {
		log.Error().Err(err)
		return "", err
	}
	log.Debug().Str("result", result).Msg("result after inserting new short url to database")
	err = g.cache.Save(ctx, result, originalUrl)
	if err != nil {
		log.Error().Err(err)
		return "", err
	}
	return result, nil
}
