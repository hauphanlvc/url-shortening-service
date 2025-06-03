package generate

import (
	"context"
	"errors"
	"log"
	"url-shortening-service/repository"
)

const MAX_RETRIES = 1000000

type GenerateService struct {
	store     repository.Store
	generator Generator
}

func NewGenerateService(store repository.Store, generator Generator) *GenerateService {
	return &GenerateService{
		store:     store,
		generator: generator,
	}
}

func (g *GenerateService) ShortUrlExist(ctx context.Context, shortUrl string) (bool, error) {
	shortUrlExist, err := g.store.CheckShortUrlExists(ctx, shortUrl)
	if err != nil {
		return false, err
	}
	return shortUrlExist, nil
}

func (g *GenerateService) InsertNewShortUrl(ctx context.Context, originalUrl string) (*string, error) {
	counter := 1
	var shortUrl string
	var err error
	for {

		shortUrl, err = g.generator.GenerateShortUrl()
		if err != nil {
			return nil, err
		}

		exists, err := g.ShortUrlExist(ctx, shortUrl)
		if err != nil {
			return nil, err
		}
		if !exists {
			break
		}

		log.Println(exists)
		log.Println(shortUrl)
		log.Printf("The short url exists, retry %d", counter)
		if counter > MAX_RETRIES {
			return nil, errors.New("cannot generate short url, always got the collision")
		}

		counter += 1
		continue
	}
	result, err := g.store.InsertNewShortUrl(ctx, originalUrl, shortUrl)
	if err != nil {
		return nil, err
	}
	return &result.ShortUrl, nil
}
