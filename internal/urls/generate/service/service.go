package generate

import (
	"context"
	"errors"
	"fmt"
	_ "log"
	"url-shortening-service/internal/cache"
	"url-shortening-service/internal/repository"

	"github.com/rs/zerolog/log"
)

const MAX_RETRIES = 1000000

var ErrCannotGernerateShortURL = errors.New("cannot generate short url, always got the collision")

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
	log.Info().Msg("Generating unique short URL...")

	for attempt := 1; attempt <= MAX_RETRIES; attempt++ {
		shortUrl, err := g.generator.GenerateShortUrl()
		if err != nil {
			return "", fmt.Errorf("failed to generate short URL: %w", err)
		}

		exists, err := g.ShortUrlExist(ctx, shortUrl)
		if err != nil {
			return "", fmt.Errorf("failed to check if short URL exists: %w", err)
		}

		if !exists {
			log.Debug().Str("shortUrl", shortUrl).Msg("Generated unique short URL.")
			return shortUrl, nil
		}

		log.Debug().Int("attempt", attempt).Str("shortUrl", shortUrl).Msg("Short URL already exists, trying again...")
	}

	log.Error().Err(ErrCannotGernerateShortURL).Msgf("failed after %d attempts", MAX_RETRIES)
	return "", ErrCannotGernerateShortURL
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
