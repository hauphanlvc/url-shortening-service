package cache

import (
	"context"
	"time"
	"url-shortening-service/internal/repository"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

const EXPIRATION_TIME_MINUTE = 1

type Cache interface {
	Save(ctx context.Context, shortUrl string, r repository.RetrieveShortUrlRow) error
	Get(ctx context.Context, shortUrl string) (*repository.RetrieveShortUrlRow, error)
}

type ShortUrlRespone struct {
	originalUrl string `redis:"originalUrl"`
	expiredAt   string `redis:"expiredAt"`
}
type DrangonFlyCache struct {
	client *redis.Client
}

func NewDrangonFlyCache(client *redis.Client) *DrangonFlyCache {
	return &DrangonFlyCache{
		client: client,
	}
}

func (d *DrangonFlyCache) Save(ctx context.Context, shortUrl string, r repository.RetrieveShortUrlRow) error {
	// err := d.client.Set(ctx, shortUrl, originalUrl, 1).Err()
	hashFields := []string{
		"originalUrl", r.OriginalUrl,
		"expiredAt", r.ExpiredAt.Format(time.RFC3339),
	}
	status, err := d.client.HSet(ctx, shortUrl, hashFields).Result()

	// if err != nil {
	// 	panic(err)
	// }

	// status := d.client.Set(ctx, shortUrl, r, time.Duration(EXPIRATION_TIME_MINUTE)*time.Minute)
	log.Debug().Msgf("Status cmd %v", status)
	return err
}

func (d *DrangonFlyCache) Get(ctx context.Context, shortUrl string) (*repository.RetrieveShortUrlRow, error) {
	log.Info().Msgf("get shortUrl from cache")

	originalUrl, err := d.client.HGet(ctx, shortUrl, "originalUrl").Result()

	if err == redis.Nil {
		log.Error().Err(err).Msgf("shortUrl %v does not exist in cache", shortUrl)
		return nil, nil
	} else if err != nil {
		log.Error().Err(err).Msgf("Cannot get the cache of shortUrl, maybe the cache server is error")
		return nil, err
	}

	expiredAt, err := d.client.HGet(ctx, shortUrl, "expiredAt").Result()
	if err != nil {
	}
	parsedExpiredAt, err := time.Parse(time.RFC3339, expiredAt)

	if err != nil {
		return nil, err
	}
	// originalUrl, err := d.client.Get(ctx, shortUrl).Result()
	log.Debug().Str("originalUrl", originalUrl)
	log.Debug().Str("expiredAt", expiredAt)
	log.Debug().Str("originalUrl", originalUrl).Msgf("get from the cache successfully")
	log.Debug().Str("expiredAt", originalUrl).Msgf("get from the cache successfully")
	return &repository.RetrieveShortUrlRow{OriginalUrl: originalUrl, ExpiredAt: parsedExpiredAt}, nil
}
