package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/jixlox0/studoto-backend/internal/config"
)

type OAuthService interface {
	GetGoogleAuthURL(state string) string
	GetGitHubAuthURL(state string) string
	ExchangeGoogleCode(code string) (*OAuthUser, error)
	ExchangeGitHubCode(code string) (*OAuthUser, error)
}

type OAuthUser struct {
	ID        string
	Email     string
	Name      string
	AvatarURL string
	Provider  string
}

type oauthService struct {
	googleClientID     string
	googleClientSecret string
	githubClientID     string
	githubClientSecret string
	redirectURL        string
}

func NewOAuthService(cfg config.OAuthConfig) OAuthService {
	return &oauthService{
		googleClientID:     cfg.Google.ClientID,
		googleClientSecret: cfg.Google.ClientSecret,
		githubClientID:     cfg.GitHub.ClientID,
		githubClientSecret: cfg.GitHub.ClientSecret,
		redirectURL:        cfg.RedirectURL,
	}
}

func (s *oauthService) GetGoogleAuthURL(state string) string {
	params := url.Values{}
	params.Set("client_id", s.googleClientID)
	params.Set("redirect_uri", s.redirectURL)
	params.Set("response_type", "code")
	params.Set("scope", "openid email profile")
	params.Set("state", state)
	params.Set("access_type", "offline")

	return fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?%s", params.Encode())
}

func (s *oauthService) GetGitHubAuthURL(state string) string {
	params := url.Values{}
	params.Set("client_id", s.githubClientID)
	params.Set("redirect_uri", s.redirectURL)
	params.Set("scope", "user:email")
	params.Set("state", state)

	return fmt.Sprintf("https://github.com/login/oauth/authorize?%s", params.Encode())
}

func (s *oauthService) ExchangeGoogleCode(code string) (*OAuthUser, error) {
	// Exchange code for token
	tokenURL := "https://oauth2.googleapis.com/token"
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", s.googleClientID)
	data.Set("client_secret", s.googleClientSecret)
	data.Set("redirect_uri", s.redirectURL)
	data.Set("grant_type", "authorization_code")

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	// Get user info
	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo"
	req, _ := http.NewRequest("GET", userInfoURL, nil)
	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	var userInfo struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &OAuthUser{
		ID:        userInfo.ID,
		Email:     userInfo.Email,
		Name:      userInfo.Name,
		AvatarURL: userInfo.Picture,
		Provider:  "google",
	}, nil
}

func (s *oauthService) ExchangeGitHubCode(code string) (*OAuthUser, error) {
	// Exchange code for token
	tokenURL := "https://github.com/login/oauth/access_token"
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", s.githubClientID)
	data.Set("client_secret", s.githubClientSecret)
	data.Set("redirect_uri", s.redirectURL)

	req, _ := http.NewRequest("POST", tokenURL, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.PostForm = data

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	// Get user info
	userInfoURL := "https://api.github.com/user"
	req, _ = http.NewRequest("GET", userInfoURL, nil)
	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err = client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var userInfo struct {
		ID     int    `json:"id"`
		Login  string `json:"login"`
		Name   string `json:"name"`
		Email  string `json:"email"`
		Avatar string `json:"avatar_url"`
	}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	// Get email if not in user info
	if userInfo.Email == "" {
		emailURL := "https://api.github.com/user/emails"
		req, _ = http.NewRequest("GET", emailURL, nil)
		req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
		req.Header.Set("Accept", "application/vnd.github.v3+json")

		resp, err = client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			var emails []struct {
				Email   string `json:"email"`
				Primary bool   `json:"primary"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&emails); err == nil {
				for _, e := range emails {
					if e.Primary {
						userInfo.Email = e.Email
						break
					}
				}
			}
		}
	}

	return &OAuthUser{
		ID:        fmt.Sprintf("%d", userInfo.ID),
		Email:     userInfo.Email,
		Name:      userInfo.Name,
		AvatarURL: userInfo.Avatar,
		Provider:  "github",
	}, nil
}
