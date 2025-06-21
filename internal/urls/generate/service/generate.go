package generate

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/rs/zerolog/log"
)

const SHORT_URL_LEN = 7

type Generator interface {
	GenerateShortUrl() (string, error)
}

type NanoIdGenerator struct{}

func NewNannoIdGenerator() *NanoIdGenerator {
	return &NanoIdGenerator{}
}

func (n *NanoIdGenerator) GenerateShortUrl() (string, error) {
	shortUrl, err := gonanoid.New(SHORT_URL_LEN)
	log.Logger.Info().Msgf("GenerateShortUrl")
	log.Logger.Debug().Msgf("shortUrl %s", shortUrl)
	return shortUrl, err
}
