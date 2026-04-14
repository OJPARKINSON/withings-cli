package auth

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/spf13/cobra"
)

var reDirectUrl = "http://localhost:8080/callback"

func SignIn(cmd *cobra.Command, args []string) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	configDirPath := filepath.Join(home, ".config")

	configDir, err := os.ReadDir(configDirPath)
	if err != nil {
		log.Fatal(err)
	}

	withingsPath := filepath.Join(configDirPath, "withings-cli.toml")
	if slices.ContainsFunc(configDir, func(dir os.DirEntry) bool { return dir.Name() == "withings-cli.toml" }) {

		withingsConfigBytes, _ := os.ReadFile(withingsPath)

		withingsConfig, _ := DecodeConfig(withingsConfigBytes)

		if withingsConfig.ExpiresAt > time.Now().Unix() {
			fmt.Println("✓ You are logged in")
			return
		}
	}

	// login the user
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

	var scope = rand.Text()

	clientId := os.Getenv("CLIENT_ID")

	fmt.Printf("https://account.withings.com/oauth2_user/authorize2?response_type=code&client_id=%s&scope=user.info,user.metrics,user.activity&redirect_uri=%s&state=%s", clientId, reDirectUrl, scope)

	select {
	case code := <-codeCh:
		srv.Shutdown(context.Background())

		res, err := http.PostForm("https://wbsapi.withings.net/v2/oauth2", url.Values{
			"action":        {"requesttoken"},
			"grant_type":    {"authorization_code"},
			"client_id":     {os.Getenv("CLIENT_ID")},
			"client_secret": {os.Getenv("CLIENT_SECRET")},
			"code":          {code},
			"redirect_uri":  {reDirectUrl},
		})
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		var result struct {
			Status int `json:"status"`
			Body   struct {
				AccessToken  string `json:"access_token"`
				RefreshToken string `json:"refresh_token"`
				UserID       string `json:"userid"`
				ExpiresIn    int64  `json:"expires_in"`
			} `json:"body"`
		}
		json.NewDecoder(res.Body).Decode(&result)

		cfg := Config{
			AccessToken:  result.Body.AccessToken,
			RefreshToken: result.Body.RefreshToken,
			UserId:       result.Body.UserID,
			ExpiresAt:    time.Now().Unix() + result.Body.ExpiresIn,
		}

		data, err := EncodeConfig(cfg)
		if err != nil {
			log.Fatal(err)
		}

		os.WriteFile(withingsPath, data, 0600)
		fmt.Println("\n ✓ Logged in successfully")

	case err := <-errCh:
		log.Fatal(err)
	case <-time.After(2 * time.Minute):
		log.Fatal("timed out waiting for auth")
	}

}
