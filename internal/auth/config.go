package auth

import (
	"bytes"

	toml "github.com/pelletier/go-toml/v2"
)

type Config struct {
	AccessToken  string `toml:"access_token"`
	RefreshToken string `toml:"refresh_token"`
	UserId       string `toml:"user_id"`
	ExpiresAt    int64  `toml:"expires_at"`
}

func DecodeConfig(tomlData []byte) (Config, error) {
	var conf Config
	err := toml.Unmarshal(tomlData, &conf)
	if err != nil {
		return conf, nil
	}

	return conf, nil
}

func EncodeConfig(tomlData Config) ([]byte, error) {
	var buf = new(bytes.Buffer)

	err := toml.NewEncoder(buf).Encode(tomlData)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
