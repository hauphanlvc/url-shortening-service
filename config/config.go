package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
	Database struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"sslmode"`
	} `mapstructure:"database"`
	Shortener struct {
		BaseURL    string `mapstructure:"base_url"`
		CodeLength int    `mapstructure:"code_length"`
	} `mapstructure:"shortener"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w", err))
	}
	viper.AutomaticEnv()

	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
		return nil, err
	}
	return &config, nil
}
