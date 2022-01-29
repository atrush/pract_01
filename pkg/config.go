package pkg

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	ServerPort      string `env:"SERVER_ADDRESS" validate:"required,hostname_port"`
	BaseURL         string `env:"BASE_URL" validate:"required,url"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" validate:"-"`
}

const (
	defServerPort  = ":8080"
	defBaseURL     = "http://localhost:8080"
	defFileStorage = ""
)

func NewConfig() (*Config, error) {
	cfg := Config{}
	cfg.readFlagConfig()
	cfg.readEnvConfig()

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("ошибка инициализации конфига: %w", err)
	}
	return &cfg, nil
}

func (c *Config) readFlagConfig() {
	flag.StringVar(&c.ServerPort, "a", defServerPort, "порт HTTP-сервера <:port>")
	flag.StringVar(&c.BaseURL, "b", defBaseURL, "базовый URL для сокращенных ссылок <http://localhost:port>")
	flag.StringVar(&c.FileStoragePath, "f", defFileStorage, "путь до файла с сокращёнными URL")

	flag.Parse()
}

func (c *Config) readEnvConfig() error {
	envConfig := &Config{}

	if err := env.Parse(envConfig); err != nil {
		return fmt.Errorf("ошибка чтения переменных окружения:%w", err)
	}

	if envConfig.BaseURL != "" {
		c.BaseURL = envConfig.BaseURL
	}
	if envConfig.ServerPort != "" {
		c.ServerPort = envConfig.ServerPort
	}
	if envConfig.FileStoragePath != "" {
		c.FileStoragePath = envConfig.FileStoragePath
	}

	return nil
}

func (c *Config) Validate() error {
	validate := validator.New()

	if err := validate.Struct(c); err != nil {
		return fmt.Errorf("ошибка валидации конфига: %w", err)
	}

	return nil
}
