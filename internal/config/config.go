package config

import (
	"encoding/json"
	"os"
)

// Config represents the server configuration
type Config struct {
	Server struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"server"`
	Redis struct {
		Address  string `json:"address"`
		Password string `json:"password"`
		DB       int    `json:"db"`
	} `json:"redis"`
	Game struct {
		MaxPlayers       int `json:"maxPlayers"`
		BattleTimeout    int `json:"battleTimeout"`
		MessageQueueSize int `json:"messageQueueSize"`
	} `json:"game"`
}

// LoadConfig loads the configuration from a JSON file
func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
