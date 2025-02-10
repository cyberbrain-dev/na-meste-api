// Contains tools for loading the configuration
package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

// Represents a configuration of the app
type Configuration struct {
	Env                string             `yaml:"env"`
	HTTPServer         HTTPServer         `yaml:"http_server"`
	PostgresConnection PostgresConnection `yaml:"postgres_connection"`
}

// Represents a config for the app's server
type HTTPServer struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

// Represents a config for connecting to the db
type PostgresConnection struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
}

// Loads a configuration
func MustLoad() Configuration {
	// loading the env variables
	godotenv.Load(".env")

	// getting the path to config
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("cannot find config path variable")
	}

	// checking if the file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("cannot find the config file")
	}

	var cfg Configuration

	// reading the config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("cannot read the config")
	}

	return cfg
}
