package cache

type DragonflyCache struct {
	client string
}

func (d *DragonflyCache) Save(shortUrl, originalUrl string) error {
	return nil
}

func (d *DragonflyCache) Get(shortUrl string) (string, error) {
	return "", nil
}
