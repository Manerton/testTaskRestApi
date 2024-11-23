package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DbCongig struct {
	DB_HOST     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
	DB_PORT     string
	SSL_MODE    string
}

type ServerConfig struct {
	APP_HOST string
	APP_PORT string
}

type Config struct {
	DbCongig
	ServerConfig
}

func (cfg *Config) GetDataSourceName() string {
	return fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.DB_HOST,
		cfg.DB_USER,
		cfg.DB_NAME,
		cfg.DB_PASSWORD,
		cfg.SSL_MODE,
	)
}

func (cfg *Config) GetAddress() string {
	return fmt.Sprintf("%s:%s", cfg.APP_HOST, cfg.APP_PORT)
}

func GetConfig(configPath string) *Config {
	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
		if configPath == "" {
			log.Fatalf("CONFIG PATH is not set")
		}
	}

	if err := godotenv.Load(configPath); err != nil {
		log.Fatalf("config file does not exist")
	}
	var myconfig Config

	// read db
	myconfig.DB_HOST = os.Getenv("DB_HOST")
	myconfig.DB_USER = os.Getenv("DB_USER")
	myconfig.DB_PASSWORD = os.Getenv("DB_PASSWORD")
	myconfig.DB_NAME = os.Getenv("DB_NAME")
	myconfig.SSL_MODE = os.Getenv("SSL_MODE")
	// read server
	myconfig.APP_HOST = os.Getenv("APP_HOST")
	myconfig.APP_PORT = os.Getenv("APP_PORT")

	return &myconfig
}
