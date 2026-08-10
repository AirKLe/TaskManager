package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Server   ServerConfig
	DataBase DataBaseConfig
}

type ServerConfig struct {
	Port string `json:"port"`
}

type DataBaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
	SSLMode  string `json:"sslmode"`
}

func (c *Config) DataBaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DataBase.User,
		c.DataBase.Password,
		c.DataBase.Host,
		c.DataBase.Port,
		c.DataBase.Name,
		c.DataBase.SSLMode,
	)
}

func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config

	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
