package generate_test

import(
	"testing"
	"url-shortening-service/internal/urls/generate/service"
	"github.com/stretchr/testify/assert"
)

func TestNanoIdGenerator_GenerateShortUrl(t *testing.T) {
	generator := generate.NewNanoIdGenerator()

	for range 10 {
		shortURL, err := generator.GenerateShortUrl()

		assert.NoError(t, err, "error should be nil when generating short URL")
		assert.NotEmpty(t, shortURL, "short URL should not be empty")
		assert.Len(t, shortURL, generate.SHORT_URL_LEN, "short URL should be of correct length")
	}
}

