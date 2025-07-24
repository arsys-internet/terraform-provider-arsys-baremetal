package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetErrorDescription(body []byte) (*string, error) {
	type errorResponse struct {
		Type    string `json:"type"`
		Message string `json:"message"`
		Errors  string `json:"errors"`
	}

	bodyReader := bytes.NewReader(body)
	var response errorResponse
	decoder := json.NewDecoder(bodyReader)

	if err := decoder.Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode error response: %w", err)
	}

	if response.Errors != "" {
		return &response.Errors, nil
	}

	if response.Message != "" {
		return &response.Message, nil
	}

	if response.Type != "" {
		return &response.Type, nil
	}

	return nil, fmt.Errorf("no error description found")
}

func HandleErrorResponse(resp *http.Response, expectedStatus int, operation string) error {
	switch resp.StatusCode {
	case http.StatusInternalServerError:
		return fmt.Errorf("internal server error during %s", operation)
	case http.StatusNotFound:
		return fmt.Errorf("resource not found")
	case http.StatusUnauthorized:
		return fmt.Errorf("unauthorized access during %s - please check your credentials", operation)
	case expectedStatus:
		return nil
	default:
		return handleUnexpectedStatus(resp, operation)
	}
}

func handleUnexpectedStatus(resp *http.Response, operation string) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s: %d", operation, resp.StatusCode)
	}

	description, err := GetErrorDescription(body)
	if err != nil {
		return fmt.Errorf("%s: %d - %s", operation, resp.StatusCode, string(body))
	}

	return fmt.Errorf("%s: %d - %s", operation, resp.StatusCode, *description)
}
