package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	EnvLocal = "local"
	EnvProd  = "prod"
)

type Config struct {
	Application
	Telegram Telegram `yaml:"telegram" env-required:"true"`
	Database Database `yaml:"database" env-required:"true"`
}

type Application struct {
	Env string `yaml:"env" env-required:"true"`
}

type Telegram struct {
	Token   string        `yaml:"token" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-default:"60s"`
}

type Database struct {
	Path string `yaml:"path" env-default:"./database/schedule.db"`
}

func MustLoad() *Config {
	configPath := getConfigPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	cfg := &Config{}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		log.Fatalf("Error reading config: %s", err)
	}

	return cfg
}

func getConfigPath() string {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/local.yaml"
	}
	return configPath
}
