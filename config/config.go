package config

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	PendingTime   time.Duration `yaml:"pending_time"`
	CheckInterval time.Duration `yaml:"check_interval"`
	HTTPServer    `yaml:"http_server"`
	Storage       `yaml:"storage"`
}

type HTTPServer struct {
	Address      string        `yaml:"address"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

func (s *HTTPServer) Server() *http.Server {
	return &http.Server{
		Addr:         s.Address,
		ReadTimeout:  s.ReadTimeout,
		WriteTimeout: s.WriteTimeout,
		IdleTimeout:  s.IdleTimeout,
	}
}

type Storage struct {
	Port         string `yaml:"port"`
	User         string `yaml:"user"`
	DatabaseName string `yaml:"db_name"`
	Host         string `yaml:"host"`
}

func (s *Storage) StorageURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		s.User,
		os.Getenv("POSTGRES_PASSWORD"),
		s.Host,
		s.Port,
		s.DatabaseName)
}

func NewConfig() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		return nil, fmt.Errorf("config path is empty")
	}

	_, err := os.Stat(configPath)

	if os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}

	var cfg Config

	err = cleanenv.ReadConfig(configPath, &cfg)

	if err != nil {
		return nil, fmt.Errorf("can not read config: %s", err)
	}

	return &cfg, nil
}
