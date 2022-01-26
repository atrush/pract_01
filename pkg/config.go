package pkg

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerPort      string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func NewConfig() *Config {
	return &Config{
		ServerPort: ":8080",
		BaseURL:    "http://localhost:8080",
	}
}

func ReadEnvConfig(cfg *Config) {
	envConfig := &Config{}
	_ = env.Parse(envConfig)
	if envConfig.BaseURL != "" {
		cfg.BaseURL = envConfig.BaseURL
	}
	if envConfig.ServerPort != "" {
		cfg.ServerPort = envConfig.ServerPort
	}
	if envConfig.FileStoragePath != "" {
		cfg.FileStoragePath = envConfig.FileStoragePath
	}
}

func ReadFlagConfig(cfg *Config) {
	serverPort := flag.String("a", cfg.ServerPort, "порт HTTP-сервера")
	baseURL := flag.String("b", cfg.BaseURL, "базовый URL для сокращенных ссылок")
	fileStoragePath := flag.String("f", cfg.FileStoragePath, "путь до файла с сокращёнными URL")
	flag.Parse()

	cfg.BaseURL = *baseURL
	cfg.FileStoragePath = *fileStoragePath
	cfg.ServerPort = *serverPort
}
