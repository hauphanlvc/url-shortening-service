package cache

type Cache interface {
	Save(shortUrl, originalUrl string) error
	Get(shortUrl string) (string, error)
}
