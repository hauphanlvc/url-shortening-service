package config_test

import(
	"testing"
	"url-shortening-service/config"
	"github.com/stretchr/testify/assert"
	"os"
)


func TestLoadConfig_Success(t *testing.T) {
	// Set up environment variables
	os.Setenv("APP_ENV", "development") // to skip godotenv.Load()
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_SSLMODE", "disable")
	os.Setenv("DB_REDIS_HOST", "localhost:6379")

	cfg, err := config.LoadConfig(".")

	assert.NoError(t, err)
	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, "5432", cfg.Port)
	assert.Equal(t, "testuser", cfg.User)
	assert.Equal(t, "testpass", cfg.Password)
	assert.Equal(t, "testdb", cfg.DBName)
	assert.Equal(t, "disable", cfg.SSLMode)
	assert.Equal(t, "localhost:6379", cfg.RedisHost)
}

func TestLoadConfig_MissingEnvVars(t *testing.T) {
	os.Clearenv() // Clear all env vars

	os.Setenv("APP_ENV", "development") // skip godotenv

	// Not setting DB_ vars to test missing envs
	cfg, err := config.LoadConfig(".")

	assert.NoError(t, err)
	assert.Empty(t, cfg.Host)
	assert.Empty(t, cfg.Port)
	assert.Empty(t, cfg.User)
	assert.Empty(t, cfg.Password)
	assert.Empty(t, cfg.DBName)
	assert.Empty(t, cfg.SSLMode)
	assert.Empty(t, cfg.RedisHost)
}

