package generate

import (
	"context"
	"url-shortening-service/repository"
)

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

func (g *GenerateService) InsertNewShortUrl(ctx context.Context, originalUrl string) (*string, error) {
	shortUrl, err := g.generator.GenerateShortUrl()
	if err != nil {
		return nil, err
	}
	result, err := g.store.InsertNewShortUrl(ctx, originalUrl, shortUrl)
	if err != nil {
		return nil, err
	}
	return &result.ShortUrl, nil
}
