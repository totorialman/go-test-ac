package config

import (
	"fmt"
	"os"
)

type PostgresConf struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func LoadConfigDB() (string, error) {
	conf := PostgresConf{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Name:     os.Getenv("POSTGRES_DB"),
	}

	if conf.Host == "" || conf.Port == "" || conf.User == "" || conf.Password == "" || conf.Name == "" {
		return "", fmt.Errorf("missing required DB env vars")
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", conf.User, conf.Password, conf.Host, conf.Port, conf.Name), nil
}
