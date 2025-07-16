package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"terraform-provider-arsys-baremetal/internal/client"
	"terraform-provider-arsys-baremetal/internal/models"
	"terraform-provider-arsys-baremetal/internal/util"
)

var _ ApiServerServiceInterface = (*ApiServerService)(nil)

type ApiServerService struct {
	client *client.APIClient
}

type ApiServerServiceInterface interface {
	GetServer(id string) (*models.ServerResponse, error)
	GetServers() ([]models.ServerResponse, error)
	CreateServer(request *models.ServerCreateRequest) (*models.ServerResponse, error)
	UpdateServer(id string, request *models.ServerUpdateRequest) (*models.ServerResponse, error)
	DeleteServer(id string) error
}

func NewApiServerService(client *client.APIClient) *ApiServerService {
	return &ApiServerService{client: client}
}

func GetServerService(m interface{}) ApiServerServiceInterface {
	if service, ok := m.(ApiServerServiceInterface); ok {
		return service
	}

	if apiClient, ok := m.(*client.APIClient); ok {
		return NewApiServerService(apiClient)
	}

	return nil
}

func (s *ApiServerService) GetServer(id string) (*models.ServerResponse, error) {
	resp, err := s.client.Get(fmt.Sprintf("/servers/%s", id))

	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("server not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var server models.ServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&server); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	return &server, nil
}

func (s *ApiServerService) GetServers() ([]models.ServerResponse, error) {
	resp, err := s.client.Get("/servers")

	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var servers []models.ServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&servers); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return servers, nil
}

func (s *ApiServerService) CreateServer(request *models.ServerCreateRequest) (*models.ServerResponse, error) {
	resp, err := s.client.Post("/servers", request)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error creating server: %s", string(body))
	}

	var createdServer models.ServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&createdServer); err != nil {
		fmt.Printf("JSON Decode Error: %v\n", err)
		return nil, err
	}

	return &createdServer, nil
}

func (s *ApiServerService) UpdateServer(id string, request *models.ServerUpdateRequest) (*models.ServerResponse, error) {
	resp, err := s.client.Put(fmt.Sprintf("/servers/%s", id), request)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error updating server: %s", string(body))
	}

	var updatedServer models.ServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&updatedServer); err != nil {
		fmt.Printf("JSON Decode Error: %v\n", err)
		return nil, err
	}

	return &updatedServer, nil
}

func (s *ApiServerService) DeleteServer(id string) error {
	resp, err := s.client.Delete(fmt.Sprintf("/servers/%s", id))
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("server with ID %s not found", id)
	}

	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error deleting server: %s", string(body))
	}

	return nil
}

func (s *ApiServerService) GetResource(id string) (util.ResourceModel, error) {
	server, err := s.GetServer(id)
	if err != nil {
		return nil, err
	}

	if server == nil {
		return nil, nil
	}

	model, diags := models.NewServerResourceModelFromRead(context.Background(), server, nil)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting to model: %v", diags)
	}

	return model, nil
}
