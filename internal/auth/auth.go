package auth

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	os.UserHomeDir()
}

func Auth(cmd *cobra.Command, args []string) {}

func LoadToken() (string, error) {
	config, err := loadConfig()
	if err != nil {
		log.Fatal("Auth failed please try login first, ", err)
	}

	// expiresTime := time.Unix(config.ExpiresAt, 0)
	// fmt.Println(expiresTime,  time.Now())

	if config.ExpiresAt < time.Now().Unix() {
		config, err := refresh(context.Background(), config.RefreshToken)
		if err != nil {
			log.Fatal("Failed to refresh please try login first, ", err)
		}

		return config.AccessToken, nil
	} else {
		return config.AccessToken, nil
	}
}
