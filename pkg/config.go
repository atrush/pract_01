package pkg

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/go-playground/validator/v10"
	"log"
	"os"
)

//  Config stores server config params.
type Config struct {
	ServerPort      string `env:"SERVER_ADDRESS" json:"server_address"  validate:"required,hostname_port"`
	BaseURL         string `env:"BASE_URL" json:"base_url" validate:"required,url"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"file_storage_path"  validate:"-"`
	DatabaseDSN     string `env:"DATABASE_DSN" json:"database_dsn"  validate:"-"`
	EnableHTTPS     bool   `env:"ENABLE_HTTPS" json:"enable_https" envDefault:"false" validate:"-"`

	Debug         bool   `env:"SHORTENER_DEBUG" json:"-" envDefault:"false" validate:"-"`
	ConfigPath    string `env:"CONFIG" json:"-" validate:"-"`
	TrustedSubnet string `env:"TRUSTED_SUBNET" json:"trusted_subnet" validate:"-"`
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

	configPath := getConfigPath()
	if len(configPath) != 0 {
		cfg.readFileConfig(configPath)
	}

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
	flagConfig := &Config{}
	flag.StringVar(&flagConfig.ServerPort, "a", defServerPort, "порт HTTP-сервера <:port>")
	flag.StringVar(&flagConfig.BaseURL, "b", defBaseURL, "базовый URL для сокращенных ссылок <http://localhost:port>")
	flag.StringVar(&flagConfig.FileStoragePath, "f", defFileStorage, "путь до файла с сокращёнными URL")
	flag.StringVar(&flagConfig.DatabaseDSN, "d", defDatabaseDSN, "строка с адресом подключения к БД")
	flag.BoolVar(&flagConfig.Debug, "debug", defDebug, "режим отладки")
	flag.BoolVar(&flagConfig.EnableHTTPS, "s", defEnableHTTPS, "включения HTTPS в веб-сервере")
	flag.StringVar(&flagConfig.ConfigPath, "c", "", "файл конфигурации")
	flag.StringVar(&flagConfig.TrustedSubnet, "t", "", "CIDR доверенной подсети")
	flag.Parse()

	c.redefineConfig(flagConfig)
}

//  redefineConfig redefines config with new config.
//  if string not empty override
//  if bool true override
func (c *Config) redefineConfig(nc *Config) {
	if nc.BaseURL != "" {
		c.BaseURL = nc.BaseURL
	}
	if nc.ServerPort != "" {
		c.ServerPort = nc.ServerPort
	}
	if nc.FileStoragePath != "" {
		c.FileStoragePath = nc.FileStoragePath
	}
	if nc.DatabaseDSN != "" {
		c.DatabaseDSN = nc.DatabaseDSN
	}
	if nc.TrustedSubnet != "" {
		c.TrustedSubnet = nc.TrustedSubnet
	}
	if nc.Debug {
		c.Debug = nc.Debug
	}
	if nc.EnableHTTPS {
		c.EnableHTTPS = nc.EnableHTTPS
	}
}

//  readEnvConfig redefines config params with environment params.
func (c *Config) readEnvConfig() error {
	envConfig := &Config{}

	if err := env.Parse(envConfig); err != nil {
		return fmt.Errorf("ошибка чтения переменных окружения:%w", err)
	}

	c.redefineConfig(envConfig)
	return nil
}

func (c *Config) readFileConfig(path string) {
	fileConfig := Must(parseFromFile(path))

	c.redefineConfig(fileConfig)
}

func Must(c *Config, err error) *Config {
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func getConfigPath() string {
	path := ""

	// tries to get config path from env
	path = os.Getenv("CONFIG")

	//  tries to get config path from flags
	args := os.Args
	flagIndex := indexFirst(args, "-c")
	if flagIndex > 0 && len(args) >= flagIndex+1 {
		path = args[flagIndex+1]
	}
	return path
}

func parseFromFile(path string) (*Config, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("error reading config from file: %w", err)
	}
	defer file.Close()

	var c Config // целевой объект
	if err := json.NewDecoder(file).Decode(&c); err != nil {
		return nil, fmt.Errorf("error reading config from file: %w", err)
	}
	return &c, nil
}

func indexFirst(list []string, el string) int {
	for i, vs := range list {
		if el == vs {
			return i
		}
	}
	return -1
}
