package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Host     string `mapstructure:"HOST"`
	Port     string `mapstructure:"PORT"`
	User     string `mapstructure:"USER"`
	Password string `mapstructure:"PASSWORD"`
	DBName   string `mapstructure:"NAME"`
	SSLMode  string `mapstructure:"SSLMODE"`
	RedisHost string `mapstructure:"REDIS_HOST"`
}

func LoadLocalEnv() error {
	env := os.Getenv("APP_ENV")
	if env == "development" {
		return nil
	}
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}
	return nil
}

func LoadConfig(path string) (*Config, error) {
	if err := LoadLocalEnv(); err != nil {
		return nil, err
	}
	viper.SetEnvPrefix("DB")
	viper.AutomaticEnv()
	envs := []string{"HOST", "PORT", "USER", "PASSWORD", "NAME", "SSLMODE", "REDIS_HOST"}
	for _, value := range envs {
		if err := viper.BindEnv(value); err != nil {
			log.Fatalf("failed to bind env variable %s: %v", value, err)
		}
	}

	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
		return nil, err
	}
	return &config, nil
}
