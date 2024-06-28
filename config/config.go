package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// struct untuk konfigurasi DB
type DbConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	Driver   string //untuk nama driver = postgres, mysql dll
}
type AppConfig struct {
	AppPort string
}

type SecurityConfig struct {
}

// struct gabungan dari semua config
type Config struct {
	DbConfig
	AppConfig
	SecurityConfig
}

func (c *Config) readConfig() error {
	//untuk membaca file .env
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	// ngecek value
	c.DbConfig = DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Driver:   os.Getenv("DB_DRIVER"),
	}
	c.AppConfig = AppConfig{
		AppPort: os.Getenv("PORT_APP"),
	}

	if c.DbConfig.Host == "" || c.DbConfig.Port == "" || c.DbConfig.Name == "" || c.DbConfig.User == "" || c.DbConfig.Password == "" || c.DbConfig.Driver == "" {
		return errors.New("environtment is empty")
	}
	return nil
}

func NewConfig() (*Config, error) {
	config := &Config{}
	if err := config.readConfig(); err != nil {
		return nil, err
	}
	return config, nil
}
