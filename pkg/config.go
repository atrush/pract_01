package pkg

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slices"
	"os"
)

//  Config stores server config params.
type Config struct {
	ServerPort      string `env:"SERVER_ADDRESS" json:"server_address"  validate:"required,hostname_port"`
	BaseURL         string `env:"BASE_URL" json:"base_url" validate:"required,url"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"file_storage_path"  validate:"-"`
	DatabaseDSN     string `env:"DATABASE_DSN" json:"database_dsn"  validate:"-"`
	EnableHTTPS     bool   `env:"ENABLE_HTTPS" json:"enable_https" envDefault:"false" validate:"-"`

	Debug      bool   `env:"SHORTENER_DEBUG" json:"-" envDefault:"false" validate:"-"`
	ConfigPath string `env:"CONFIG" json:"-" validate:"-"`
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
	cfg.tryReadFileConfig()
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
	flag.BoolVar(&c.EnableHTTPS, "c", defEnableHTTPS, "файл конфигурации")
	flag.Parse()
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

//  tryReadFromFile if exist flag -c or env CONFIG tries to read config from path.
//  If read ok, overrides config with new config
func (c *Config) tryReadFileConfig() {
	path := getConfigPath()
	if len(path) == 0 {
		return
	}

	fileConfig, err := readConfigFromFile(path)
	if err != nil {
		return
	}

	c.redefineConfig(fileConfig)
}

func getConfigPath() string {
	path := ""

	// tries to get config path from env
	path = os.Getenv("CONFIG")

	//  tries to get config path from flags
	args := os.Args
	flagIndex := slices.Index(args, "-c")
	if flagIndex > 0 && len(args) >= flagIndex+1 {
		path = args[flagIndex+1]
	}
	return path
}

func readConfigFromFile(path string) (*Config, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0777)
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
