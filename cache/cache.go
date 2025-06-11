package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

const EXPIRATION_TIME_MINUTE = 1

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
	// err := d.client.Set(ctx, shortUrl, originalUrl, 1).Err()
	status := d.client.Set(ctx, shortUrl, originalUrl, time.Duration(EXPIRATION_TIME_MINUTE)*time.Minute)
	log.Debug().Msgf("Status cmd %v", status)
	return status.Err()
}

func (d *DrangonFlyCache) Get(ctx context.Context, shortUrl string) (string, error) {
	log.Info().Msgf("get shortUrl from cache")
	originalUrl, err := d.client.Get(ctx, shortUrl).Result()
	log.Debug().Str("originalUrl", originalUrl)
	if err == redis.Nil {
		log.Error().Err(err).Msgf("shortUrl %v does not exist in cache", shortUrl)
		return "", nil
	} else if err != nil {
		log.Error().Err(err).Msgf("Cannot get the cache of shortUrl, maybe the cache server is error")
		return "", err
	}
	log.Debug().Str("originalUrl", originalUrl).Msgf("get from the cache successfully")
	return originalUrl, nil
}
