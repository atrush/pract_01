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

const (
	DEF_SERVER_PORT       = ":8080"
	DEF_BASE_URL          = "http://localhost:8080"
	DEF_FILE_STORAGE_PATH = ""
)

func NewConfig() *Config {
	cfg := Config{}
	cfg.readFlagConfig()
	cfg.readEnvConfig()
	return &cfg
}

func (c *Config) readFlagConfig() {
	serverPort := flag.String("a", DEF_SERVER_PORT, "порт HTTP-сервера")
	baseURL := flag.String("b", DEF_BASE_URL, "базовый URL для сокращенных ссылок")
	fileStoragePath := flag.String("f", DEF_FILE_STORAGE_PATH, "путь до файла с сокращёнными URL")
	flag.Parse()

	c.BaseURL = *baseURL
	c.FileStoragePath = *fileStoragePath
	c.ServerPort = *serverPort
}

func (c *Config) readEnvConfig() {
	envConfig := &Config{}
	_ = env.Parse(envConfig)
	if envConfig.BaseURL != "" {
		c.BaseURL = envConfig.BaseURL
	}
	if envConfig.ServerPort != "" {
		c.ServerPort = envConfig.ServerPort
	}
	if envConfig.FileStoragePath != "" {
		c.FileStoragePath = envConfig.FileStoragePath
	}
}
