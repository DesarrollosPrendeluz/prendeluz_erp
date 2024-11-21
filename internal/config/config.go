package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	Port int `yaml:"port"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
}

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

func LoadConfig(env string) (*Config, error) {
	filename := fmt.Sprintf("internal/config/%s.yaml", env)
	data, err := os.ReadFile(filename)

	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error unmarshalling config file: %v", err)
	}

	if port := os.Getenv("SERVER_PORT"); port != "" {
		fmt.Println(port)
		config.Server.Port, _ = strconv.Atoi(port)
	}
	if host := os.Getenv("DB_HOST"); host != "" {
		config.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		config.Database.Port, _ = strconv.Atoi(port)
	}
	if username := os.Getenv("DB_USERNAME"); username != "" {
		config.Database.Username = username
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		config.Database.Password = password
	}
	if dbname := os.Getenv("DB_NAME"); dbname != "" {
		config.Database.Dbname = dbname
	}

	return &config, nil
}
