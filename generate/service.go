package generate

import (
	"context"
	"url-shortening-service/repository"
)

type GenerateService struct {
	store repository.Store
}

func NewGenerateService(store repository.Store) *GenerateService {
	return &GenerateService{
		store: store,
	}
}

func (g *GenerateService) InsertNewShortUrl(ctx context.Context, originalUrl string) (*string, error) {
	shortUrl, err := GenerateShortUrl(originalUrl)
	if err != nil {
		return nil, err
	}
	result, err := g.store.InsertNewShortUrl(ctx, originalUrl, shortUrl)
	if err != nil {
		return nil, err
	}
	return &result.ShortUrl, nil
}
