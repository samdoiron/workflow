package workflow

import (
	"io/ioutil"
	"log"

	"github.com/BurntSushi/toml"
)

// DatabaseConfig contains the global configuration for the database
type DatabaseConfig struct {
	Name string
}

// ServerConfig contains global configuration for the HTTP(s) server
type ServerConfig struct {
	Port int
}

// Config is the global application configuration
type Config struct {
	DB     DatabaseConfig `toml:"database"`
	Server ServerConfig
}

const configPath = "./config.toml"

// LoadConfig loads the global application configuration from disk
func LoadConfig() Config {
	r, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	var config Config
	if err = toml.Unmarshal(r, &config); err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	return config
}
