package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Database struct {
		Host     string `mapstructure:"POSTGRES_HOST"`
		Port     int    `mapstructure:"POSTGRES_PORT"`
		User     string `mapstructure:"POSTGRES_USER"`
		Password string `mapstructure:"POSTGRES_PASSWORD"`
		DBName   string `mapstructure:"POSTGRES_DB"`
	} `mapstructure:"database"`
	Service struct {
		ServiceName string `mapstructure:"SERVICE_NAME"`
		Port        int    `mapstructure:"SERVICE_PORT"`
	} `mapstructure:"service"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	viper.AutomaticEnv()

	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
		return nil, err
	}
	return &config, nil
}
