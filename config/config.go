package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	Driver   string
}

type ApiConfig struct {
	ApiPort string
}

type Config struct {
	DBConfig
	ApiConfig
}

func (c *Config) readConfig() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("missing env file %v", err.Error())
	}
	c.DBConfig = DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		Driver:   os.Getenv("DB_DRIVER"),
	}

	c.ApiConfig = ApiConfig{ApiPort: os.Getenv("API_PORT")}


	if c.Host == "" || c.Port == "" || c.User == "" || c.Name == "" || c.Driver == "" || c.ApiPort == "" {
		return fmt.Errorf("missing required environment")
	}

	return nil

}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := cfg.readConfig(); err != nil {
		return nil, err
	}
	return cfg, nil
}
