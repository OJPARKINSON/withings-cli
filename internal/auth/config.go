package auth

import (
	"bytes"
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
	lastUpdated int64 `toml:"last_updated"`
}

func configPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	var configDirPath = filepath.Join(home, ".config", "withings")
	os.MkdirAll(configDirPath, 0755)
	var withingsPath = filepath.Join(configDirPath, "withings-cli.toml")

	return withingsPath
}

func UpdateLastUpdated(lastUpdated int64){
	config, err := loadConfig()
	if err != nil {
		log.Panic(err)
	}

	config.lastUpdated = lastUpdated

	writeConfig(config)
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
