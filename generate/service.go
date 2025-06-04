package generate

import (
	"context"
	"errors"
	"log"
)

const MAX_RETRIES = 1000000

type Store interface {
	InsertNewShortUrl(ctx context.Context, originalUrl, shortUrl string) (string, error)
	RertrieveShortUrl(ctx context.Context, shortUrl string) (string, error)
	CheckShortUrlExists(ctx context.Context, shortUrl string) (bool, error)
}

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
	shortUrlExist, err := g.store.CheckShortUrlExists(ctx, shortUrl)
	if err != nil {
		return false, err
	}
	return shortUrlExist, nil
}

func (g *GenerateService) GetUniqueShortUrl(ctx context.Context) (string, error) {
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

		log.Println(exists)
		log.Println(shortUrl)
		log.Printf("The short url exists, retry %d", counter)
		if counter > MAX_RETRIES {
			return "", errors.New("cannot generate short url, always got the collision")
		}

		counter += 1
		continue
	}
	return shortUrl, nil
}

func (g *GenerateService) InsertNewShortUrl(ctx context.Context, originalUrl string) (string, error) {
	shortUrl, err := g.GetUniqueShortUrl(ctx)
	if err != nil {
		return "", err
	}
	result, err := g.store.InsertNewShortUrl(ctx, originalUrl, shortUrl)
	if err != nil {
		return "", err
	}
	err = g.cache.Save(ctx, result, originalUrl)
	if err != nil {
		return "", err
	}
	return result, nil
}
