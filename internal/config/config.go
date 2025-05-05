package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	configFileName = ".gatorconfig.json"
)

type Config struct {
	DBURL string `json:"db_url"`
	User  string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, configFileName), nil
}

func write(cfg Config) error {
	content, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	confFile, err := getConfigFilePath()
	if err != nil {
		return err
	}
	if err := os.WriteFile(confFile, content, 0o666); err != nil {
		return err
	}
	return nil
}

func Read() (Config, error) {
	confFile, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	content, err := os.ReadFile(confFile)
	if err != nil {
		return Config{}, err
	}
	var config Config
	if err := json.Unmarshal(content, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}

func (c *Config) SetUser(userName string) error {
	c.User = userName
	if err := write(*c); err != nil {
		return err
	}
	return nil
}
