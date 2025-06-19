package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type APIClient struct {
	HTTPClient *http.Client
	BaseURL    string
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
	}
}

func (c *APIClient) sendRequest(method, path string, body interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, path)

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return c.HTTPClient.Do(req)
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
