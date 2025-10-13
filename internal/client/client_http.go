package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type APIClient struct {
	HTTPClient *http.Client
	BaseURL    string
	MaxRetries int
	BaseDelay  time.Duration
	MaxDelay   time.Duration
}

func NewAPIClient(apiToken string, url string) *APIClient {
	client := &http.Client{
		Timeout: time.Second * 30,

		Transport: &http.Transport{
			TLSHandshakeTimeout: 30 * time.Second,
			IdleConnTimeout:     90 * time.Second,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 25,

			TLSClientConfig: &tls.Config{
				MinVersion:         tls.VersionTLS12,
				InsecureSkipVerify: false,
			},
		},
	}

	client.Transport = &TokenTransport{
		Token: apiToken,
		Base:  http.DefaultTransport,
	}

	return &APIClient{
		HTTPClient: client,
		BaseURL:    url,
		MaxRetries: 10,
		BaseDelay:  time.Second * 6,
		MaxDelay:   time.Second * 30,
	}
}

func (c *APIClient) calculateDelay(attempt int, rateLimitReset string) time.Duration {
	if rateLimitReset != "" {
		if resetTime, err := strconv.ParseInt(rateLimitReset, 10, 64); err == nil {
			resetTimestamp := time.Unix(resetTime, 0)
			now := time.Now()
			delay := resetTimestamp.Sub(now)

			if delay > 0 && delay <= c.MaxDelay {
				return delay
			}
		}
	}

	delay := time.Duration(float64(c.BaseDelay) * math.Pow(2, float64(attempt)))
	if delay > c.MaxDelay {
		delay = c.MaxDelay
	}

	return delay
}

func (c *APIClient) sendRequest(method, path string, body interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, path)

	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	for attempt := 0; attempt <= c.MaxRetries; attempt++ {
		var bodyReader io.Reader
		if bodyBytes != nil {
			bodyReader = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequest(method, url, bodyReader)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}

		if resp.StatusCode != http.StatusTooManyRequests {
			return resp, nil
		}

		_ = resp.Body.Close()

		if attempt >= c.MaxRetries {
			return nil, err
		}

		delay := c.calculateDelay(attempt, resp.Header.Get("X-Rate-Limit-Reset"))
		tflog.Warn(context.Background(), "Rate limited (429), waiting before retry", map[string]interface{}{
			"delay":         delay.String(),
			"current_retry": attempt + 2,
			"max_retries":   c.MaxRetries + 1,
		})
		time.Sleep(delay)
	}

	return nil, fmt.Errorf("unexpected error in retry logic")
}

func (c *APIClient) Get(path string) (*http.Response, error) {
	return c.sendRequest("GET", path, nil)
}

func (c *APIClient) Post(path string, body interface{}) (*http.Response, error) {
	return c.sendRequest("POST", path, body)
}

func (c *APIClient) Put(path string, body interface{}) (*http.Response, error) {
	return c.sendRequest("PUT", path, body)
}

func (c *APIClient) Patch(path string, body interface{}) (*http.Response, error) {
	return c.sendRequest("PATCH", path, body)
}

func (c *APIClient) Delete(path string) (*http.Response, error) {
	return c.sendRequest("DELETE", path, nil)
}
