package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Save(ctx context.Context, shortUrl, originalUrl string) error
	Get(ctx context.Context, shortUrl string) (string, error)
}

type DrangonFlyCache struct {
	client *redis.Client
}

func NewDrangonFlyCache(client *redis.Client) *DrangonFlyCache {
	return &DrangonFlyCache{
		client: client,
	}
}

func (d *DrangonFlyCache) Save(ctx context.Context, shortUrl, originalUrl string) error {
	err := d.client.Set(ctx, shortUrl, originalUrl, 1).Err()
	return err
}

func (d *DrangonFlyCache) Get(ctx context.Context, shortUrl string) (string, error) {
	originalUrl, err := d.client.Get(ctx, shortUrl).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return originalUrl, nil
}
