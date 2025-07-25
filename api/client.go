package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
	"crowdstrike-cli/config"
)

// Client represents a CrowdStrike API client
type Client struct {
	BaseURL      string
	ClientID     string
	ClientSecret string
	AccessToken  string
	HTTPClient   *http.Client
}

// TokenResponse represents the OAuth token response
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// NewClient creates a new CrowdStrike API client
func NewClient(cfg *config.Config) *Client {
	return &Client{
		BaseURL:      cfg.BaseURL,
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Authenticate gets an access token from CrowdStrike
func (c *Client) Authenticate() error {
	tokenURL := c.BaseURL + "/oauth2/token"

	data := url.Values{}
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("authentication failed: %s", string(body))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	c.AccessToken = tokenResp.AccessToken
	return nil
}

// makeRequest makes an authenticated request to the CrowdStrike API
func (c *Client) makeRequest(method, endpoint string) ([]byte, error) {
	if c.AccessToken == "" {
		if err := c.Authenticate(); err != nil {
			return nil, err
		}
	}

	url := c.BaseURL + endpoint

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		// Token might be expired, try to re-authenticate
		if err := c.Authenticate(); err != nil {
			return nil, err
		}
		// Retry the request with new token
		return c.makeRequest(method, endpoint)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// GetDetectionSummary gets detection summary metrics
func (c *Client) GetDetectionSummary() ([]byte, error) {
}

// GetHostSummary gets host summary metrics
func (c *Client) GetHostSummary() ([]byte, error) {
}

// GetIncidentSummary gets incident summary metrics
func (c *Client) GetIncidentSummary() ([]byte, error) {
}
