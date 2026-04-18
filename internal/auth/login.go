package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var redirectURL = "http://localhost:8080/callback"

func SignIn(cmd *cobra.Command, args []string) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	configDirPath := filepath.Join(home, ".config", "withings")

	withingsPath := filepath.Join(configDirPath, "withings-cli.toml")
	if withingsConfigBytes, err := os.ReadFile(withingsPath); err == nil {

		fmt.Println(withingsConfigBytes)
		withingsConfig, err := decodeConfig(withingsConfigBytes)
		if err != nil {
			log.Fatal(err)
		}

		if withingsConfig.ExpiresAt > time.Now().Unix() {
			fmt.Println("✓ You are logged in")
			return
		}
	}

	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)

	srv := &http.Server{Addr: ":8080"}
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errCh <- fmt.Errorf("no code in callback")
			fmt.Fprintf(w, "Error: no code received. Close this tab.")
			return
		}
		fmt.Fprintf(w, "✓ Authenticated! You can close this tab.")
		codeCh <- code
	})

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	var state = rand.Text()

	clientId := os.Getenv("CLIENT_ID")

	fmt.Printf("https://account.withings.com/oauth2_user/authorize2?response_type=code&client_id=%s&scope=user.info,user.metrics,user.activity&redirect_uri=%s&state=%s", clientId, redirectURL, state)

	select {
	case code := <-codeCh:
		srv.Shutdown(context.Background())

		_, err := ExchangeCode(context.Background(), code)
		if err != nil {
			fmt.Println("Failed to login ")
			fmt.Println(err)
		}

		fmt.Println("\n ✓ Logged in successfully")

	case err := <-errCh:
		log.Fatal(err)
	case <-time.After(2 * time.Minute):
		log.Fatal("timed out waiting for auth")
	}

}
