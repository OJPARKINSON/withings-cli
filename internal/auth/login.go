package auth

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"

	"github.com/spf13/cobra"
)

func SignIn(cmd *cobra.Command, args []string) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	configDirPath := filepath.Join(home, ".config")

	fmt.Println("Konfigurationsverzeichnis:", configDirPath)

	configDir, err := os.ReadDir(configDirPath)
	if err != nil {
		log.Fatal(err)
	}

	withingsPath := filepath.Join(configDirPath, "withings-cli.toml")
	if slices.ContainsFunc(configDir, func(dir os.DirEntry) bool { return dir.Name() == "withings-cli.toml" }) {

		withingsConfig, _ := os.ReadFile(withingsPath)

		DecodeConfig(string(withingsConfig))
	} else {

		tomlBuffer, err := EncodeConfig(Config{
			AccessToken:  "",
			RefreshToken: "",
			age:          "",
			scope:        "",
		})

		err = os.WriteFile(withingsPath, tomlBuffer, 0644)
		if err != nil {
			log.Fatal(err)
		}

		DecodeConfig(string(withingsConfig))
	}

}
