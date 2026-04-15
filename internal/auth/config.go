package auth

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"

	toml "github.com/pelletier/go-toml/v2"
)

type Config struct {
	AccessToken  string `toml:"access_token"`
	RefreshToken string `toml:"refresh_token"`
	UserId       string `toml:"user_id"`
	ExpiresAt    int64  `toml:"expires_at"`
}

func configPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	
	var configDirPath = filepath.Join(home, ".config")
	var withingsPath = filepath.Join(configDirPath, "withings/withings-cli.toml")

	return withingsPath
}

func writeConfig(cfg *Config) {
	data, err := EncodeConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}

	withingsPath := configPath()

	err = os.WriteFile(withingsPath, data, 0600)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("oath ")
}

func loadConfig() (*Config, error) {
	withingsPath := configPath()

	withingsConfigBytes, err := os.ReadFile(withingsPath)
	if err != nil {
		return nil, err
	}

	withingsConfig, err := decodeConfig(withingsConfigBytes)
	if err != nil {
		return nil, err
	}

	return withingsConfig, nil
}

func decodeConfig(tomlData []byte) (*Config, error) {
	var conf Config
	err := toml.Unmarshal(tomlData, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

func EncodeConfig(tomlData *Config) ([]byte, error) {
	var buf = new(bytes.Buffer)

	err := toml.NewEncoder(buf).Encode(tomlData)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
