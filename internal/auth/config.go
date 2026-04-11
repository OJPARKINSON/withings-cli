package auth

import (
	"bytes"

	"github.com/BurntSushi/toml"
)

type Config struct {
	AccessToken  string
	RefreshToken string
	age          string
	scope        string
}

func DecodeConfig(tomlData string) (toml.MetaData, error) {
	var conf Config
	return toml.Decode(tomlData, &conf)
}

func EncodeConfig(tomlData Config) ([]byte, error) {
	var buf = new(bytes.Buffer)

	err := toml.NewEncoder(buf).Encode(tomlData)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
