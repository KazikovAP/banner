package config

import (
	logerr "banner/internal/lib/logger/logerr"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string         `yaml:"env"`
	Server   ServerConfig   `yaml:"server"`
	Postgres PostgresConfig `yaml:"postgres"`
	Jwt      JwtConfig      `yaml:"jwt"`
}

type ServerConfig struct {
	Host        string        `yaml:"host"`
	Port        string        `yaml:"port"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type JwtConfig struct {
	Secret string `yaml:"secret"`
}

const (
	localPathToConfig = "/config/config.yaml"
)

func LoadConfig() (*Config, error) {
	log.Println("Start configuration setup")

	pwdPath, err := os.Getwd()
	if err != nil {
		log.Fatal("Can't get working directory path")
	}

	configPath := pwdPath + localPathToConfig

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("Config file does not exist ", logerr.Err(err))
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("Cannot read config")
	}

	log.Println("Configuration complete")

	return &cfg, nil
}
