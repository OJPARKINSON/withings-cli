package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	AuthURL  = "https://account.withings.com/oauth2_user/authorize2"
	TokenURL = "https://wbsapi.withings.net/v2/oauth2"

	RedirectURL = "http://localhost:8080/callback"

	Scopes = "user.info,user.metrics,user.activity"
)

func ExchangeCode(ctx context.Context, code string) (*Config, error) {
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	params := url.Values{
		"action":        {"requesttoken"},
		"grant_type":    {"authorization_code"},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"code":          {code},
		"redirect_uri":  {RedirectURL},
	}
	return doTokenRequest(ctx, params)
}

func refresh(ctx context.Context, refreshToken string) (*Config, error) {
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	params := url.Values{
		"action":        {"requesttoken"},
		"grant_type":    {"refresh_token"},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"refresh_token": {refreshToken},
	}

	return doTokenRequest(ctx, params)
}

func doTokenRequest(ctx context.Context, params url.Values) (*Config, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, TokenURL, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	var wResp WithingsTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&wResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	if wResp.Status != 0 {
		return nil, fmt.Errorf("withings API error: status=%d error=%s", wResp.Status, wResp.Error)
	}

	config := &Config{
		UserId:      wResp.Body.UserID.String(),
		AccessToken:  wResp.Body.AccessToken,
		RefreshToken: wResp.Body.RefreshToken,
		ExpiresAt:    time.Now().Unix() + wResp.Body.ExpiresIn,
	}

	writeConfig(config)

	return config, nil
}

type WithingsTokenResponse struct {
	Status int               `json:"status"`
	Body   WithingsTokenBody `json:"body"`
	Error  string            `json:"error,omitempty"`
}

type WithingsTokenBody struct {
	UserID       json.Number  `json:"userid"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}
