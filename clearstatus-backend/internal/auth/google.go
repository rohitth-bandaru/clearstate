package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// GoogleUserInfo from tokeninfo or userinfo
type GoogleUserInfo struct {
	Sub           string `json:"sub"`            // Google ID
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

// VerifyGoogleIDToken verifies the ID token with Google's tokeninfo endpoint.
// clientID is the Google OAuth client ID (same as frontend) to validate audience.
func VerifyGoogleIDToken(ctx context.Context, idToken, clientID string) (*GoogleUserInfo, error) {
	if idToken == "" {
		return nil, fmt.Errorf("id_token is required")
	}
	u := "https://oauth2.googleapis.com/tokeninfo?id_token=" + url.QueryEscape(idToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google tokeninfo returned %d", resp.StatusCode)
	}
	var payload struct {
		Audience  string `json:"aud"`
		Subject   string `json:"sub"`
		Email     string `json:"email"`
		EmailVerified string `json:"email_verified"`
		Name      string `json:"name"`
		Picture   string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	if clientID != "" && payload.Audience != clientID {
		return nil, fmt.Errorf("token audience does not match client ID")
	}
	if payload.Email == "" {
		return nil, fmt.Errorf("email not in token")
	}
	info := &GoogleUserInfo{
		Sub:           payload.Subject,
		Email:         payload.Email,
		EmailVerified: strings.ToLower(payload.EmailVerified) == "true",
		Name:          payload.Name,
		Picture:       payload.Picture,
	}
	return info, nil
}
