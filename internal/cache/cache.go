package cache

import (
	"context"
	"fmt"
	"time"
	"url-shortening-service/internal/repository"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

const EXPIRATION_TIME_MINUTE = 5

type Cache interface {
	Save(ctx context.Context, shortUrl string, data repository.RetrieveShortUrlRow) error
	Get(ctx context.Context, shortUrl string) (*repository.RetrieveShortUrlRow, error)
	Delete(ctx context.Context, shortUrl string) error
}

type ShortUrlRespone struct {
	originalUrl string `redis:"originalUrl"`
	expiredAt   string `redis:"expiredAt"`
}
type DrangonFlyCache struct {
	client *redis.Client
}

const PREFIX_SHORT_URL = "shortUrl:"

func NewDrangonFlyCache(client *redis.Client) *DrangonFlyCache {
	return &DrangonFlyCache{
		client: client,
	}
}

func (d *DrangonFlyCache) Save(ctx context.Context, shortUrl string, data repository.RetrieveShortUrlRow) error {
	key := PREFIX_SHORT_URL + shortUrl

	fields := map[string]any{
		"originalUrl": data.OriginalUrl,
		"expiredAt":   data.ExpiredAt.Format(time.RFC3339),
	}

	// 1️⃣ Save the hash
	if _, err := d.client.HSet(ctx, key, fields).Result(); err != nil {
		return fmt.Errorf("failed to HSet cache for key %q: %w", key, err)
	}

	// 2️⃣ Set TTL
	ttl := time.Duration(EXPIRATION_TIME_MINUTE) * time.Minute
	if err := d.client.Expire(ctx, key, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set expiration for key %q: %w", key, err)
	}

	log.Debug().
		Str("key", key).
		Str("originalUrl", data.OriginalUrl).
		Time("expiredAt", data.ExpiredAt).
		Dur("ttl", ttl).
		Msg("cache saved successfully")

	return nil
}

func (d *DrangonFlyCache) Get(ctx context.Context, shortUrl string) (*repository.RetrieveShortUrlRow, error) {
	key := PREFIX_SHORT_URL + shortUrl

	originalUrl, err := d.client.HGet(ctx, key, "originalUrl").Result()
	if err == redis.Nil {
		log.Warn().Str("key", key).Msg("originalUrl does not exist in cache")
		return nil, nil
	}
	if err != nil {
		log.Error().Err(err).Str("key", key).Msg("failed to get originalUrl from cache")
		return nil, fmt.Errorf("failed to get originalUrl from cache: %w", err)
	}

	expiredAtStr, err := d.client.HGet(ctx, key, "expiredAt").Result()
	if err == redis.Nil {
		log.Warn().Str("key", key).Msg("expiredAt does not exist in cache")
		return nil, nil
	}
	if err != nil {
		log.Error().Err(err).Str("key", key).Msg("failed to get expiredAt from cache")
		return nil, fmt.Errorf("failed to get expiredAt from cache: %w", err)
	}

	expiredAt, err := time.Parse(time.RFC3339, expiredAtStr)
	if err != nil {
		log.Error().Err(err).Str("expiredAt", expiredAtStr).Msg("failed to parse expiredAt")
		return nil, fmt.Errorf("failed to parse expiredAt: %w", err)
	}

	log.Debug().
		Str("key", key).
		Str("originalUrl", originalUrl).
		Time("expiredAt", expiredAt).
		Msg("fetched data from cache successfully")

	return &repository.RetrieveShortUrlRow{
		OriginalUrl: originalUrl,
		ExpiredAt:   expiredAt,
	}, nil
}

func (d *DrangonFlyCache) Delete(ctx context.Context, shortUrl string) error {

	key := PREFIX_SHORT_URL + shortUrl
	_, err := d.client.Del(ctx, key).Result()
	if err == redis.Nil {
		log.Warn().Str("key", key).Msg("originalUrl does not exist in cache")
		return nil
	} else if err != nil {
		log.Error().Err(err).Str("key", key).Msg("failed to delete key from database")
		return fmt.Errorf("failed to delete key from cache")
	}
	return nil
}
