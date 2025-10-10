package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"
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
			TLSHandshakeTimeout: 15 * time.Second,
			IdleConnTimeout:     90 * time.Second,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 20,

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
		MaxRetries: 3,
		BaseDelay:  time.Second * 2,
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
			return nil, err
		}
	}

	for attempt := 0; attempt <= c.MaxRetries; attempt++ {
		var bodyReader io.Reader
		if bodyBytes != nil {
			bodyReader = bytes.NewBuffer(bodyBytes)
		}

		req, err := http.NewRequest(method, url, bodyReader)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			if attempt < c.MaxRetries {
				delay := c.calculateDelay(attempt, "")
				time.Sleep(delay)
				continue
			}
			return nil, err
		}

		if resp.StatusCode != http.StatusTooManyRequests {
			return resp, nil
		}

		bodyError := resp.Body.Close()
		if bodyError != nil {
			return nil, bodyError
		}

		if attempt >= c.MaxRetries {
			return resp, nil
		}

		rateLimitReset := resp.Header.Get("X-Rate-Limit-Reset")
		delay := c.calculateDelay(attempt, rateLimitReset)

		fmt.Printf("Rate limited (429), retrying in %v... (attempt %d/%d)\n",
			delay, attempt+1, c.MaxRetries)

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
