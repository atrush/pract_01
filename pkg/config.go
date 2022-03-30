package pkg

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/go-playground/validator/v10"
)

//  Config stores server config params.
type Config struct {
	ServerPort      string `env:"SERVER_ADDRESS" validate:"required,hostname_port"`
	BaseURL         string `env:"BASE_URL" validate:"required,url"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" validate:"-"`
	DatabaseDSN     string `env:"DATABASE_DSN" validate:"-"`
	Debug           bool   `env:"SHORTENER_DEBUG" envDefault:"false" validate:"-"`
	EnableHTTPS     bool   `env:"ENABLE_HTTPS" envDefault:"false" validate:"-"`
}

//  Default config params.
const (
	defServerPort  = ":8080"
	defBaseURL     = "http://localhost:8080"
	defFileStorage = ""
	defDatabaseDSN = ""
	defDebug       = false
	defEnableHTTPS = false
)

//  NewConfig inits new config.
//  Reads flag params over default params, then redefines  with environment params.
func NewConfig() (*Config, error) {
	cfg := Config{}
	cfg.readFlagConfig()
	if err := cfg.readEnvConfig(); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("ошибка инициализации конфига: %w", err)
	}
	return &cfg, nil
}

//  Validate validates config params.
func (c *Config) Validate() error {
	validate := validator.New()

	if err := validate.Struct(c); err != nil {
		return fmt.Errorf("ошибка валидации конфига: %w", err)
	}
	return nil
}

//  readFlagConfig reads flag params over default params.
func (c *Config) readFlagConfig() {
	flag.StringVar(&c.ServerPort, "a", defServerPort, "порт HTTP-сервера <:port>")
	flag.StringVar(&c.BaseURL, "b", defBaseURL, "базовый URL для сокращенных ссылок <http://localhost:port>")
	flag.StringVar(&c.FileStoragePath, "f", defFileStorage, "путь до файла с сокращёнными URL")
	flag.StringVar(&c.DatabaseDSN, "d", defDatabaseDSN, "строка с адресом подключения к БД")
	flag.BoolVar(&c.Debug, "debug", defDebug, "режим отладки")
	flag.BoolVar(&c.EnableHTTPS, "s", defEnableHTTPS, "включения HTTPS в веб-сервере")
	flag.Parse()
}

//  readEnvConfig redefines config params with environment params.
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
	if envConfig.DatabaseDSN != "" {
		c.DatabaseDSN = envConfig.DatabaseDSN
	}
	if envConfig.Debug {
		c.Debug = envConfig.Debug
	}
	if envConfig.EnableHTTPS {
		c.EnableHTTPS = envConfig.EnableHTTPS
	}
	return nil
}
