package generate

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"

	"github.com/jxskiss/base62"
)

const SHORT_URL_LEN = 7

type Hasher interface {
	Hash(url string) ([]byte, error)
}

type Md5HashMethod struct{}

func (m *Md5HashMethod) Hash(url string) ([]byte, error) {
	h := md5.New()
	io.WriteString(h, url)
	hashUrl := h.Sum(nil)
	hexHashUrl := fmt.Sprintf("%x", hashUrl[:SHORT_URL_LEN])
	decodedHashUrl, err := hex.DecodeString(hexHashUrl)
	if err != nil {
		return nil, err
	}
	return decodedHashUrl, nil
}

type Encoder interface {
	Encode(url string) (string, error)
}

type Base62EncodeMethod struct {
	hasher Hasher
}

func (b *Base62EncodeMethod) Encode(url string) ([]byte, error) {
	hashUrl, err := b.hasher.Hash(url)
	if err != nil {
		return nil, err
	}
	log.Printf("hashUrl %x %d", hashUrl, len(hashUrl))
	return base62.Encode(hashUrl)[:SHORT_URL_LEN], nil
}

func GenerateShortUrl(url string) (string, error) {
	md5hashMethod := &Md5HashMethod{}
	base62EncodeMethod := &Base62EncodeMethod{hasher: md5hashMethod}
	shortUrl, err := base62EncodeMethod.Encode(url)
	log.Printf("shortUrl %s %d", shortUrl, len(shortUrl))
	if err != nil {
		return "", err
	}
	return string(shortUrl[:]), nil
}
