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
	"terraform-provider-arsys-baremetal/internal/util/helper"
)

var _ ApiServerServiceInterface = (*ApiServerService)(nil)

type ApiServerService struct {
	client *client.APIClient
}

type ApiServerServiceInterface interface {
	GetServer(id string) (*models.ServerDetailResponse, error)
	GetServers() ([]models.ServerListResponse, error)
	CreateServer(request *models.ServerCreateRequest) (*models.ServerBaseResponse, error)
	UpdateServer(id string, request *models.ServerUpdateRequest) (*models.ServerBaseResponse, error)
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

func (s *ApiServerService) GetServer(id string) (*models.ServerDetailResponse, error) {
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

	errorResponse := helper.HandleErrorResponse(resp, http.StatusOK, "get server")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var server models.ServerDetailResponse
	if err := json.NewDecoder(resp.Body).Decode(&server); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if server.ServerType != "baremetal" {
		return nil, fmt.Errorf("server not found")
	}

	return &server, nil
}

func (s *ApiServerService) GetServers() ([]models.ServerListResponse, error) {
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

	errorResponse := helper.HandleErrorResponse(resp, http.StatusOK, "get servers")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var servers []models.ServerListResponse
	if err := json.NewDecoder(resp.Body).Decode(&servers); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	var baremetalServers []models.ServerListResponse
	for _, server := range servers {
		if server.ServerType == "baremetal" {
			baremetalServers = append(baremetalServers, server)
		}
	}

	return baremetalServers, nil
}

func (s *ApiServerService) CreateServer(request *models.ServerCreateRequest) (*models.ServerBaseResponse, error) {
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

	errorResponse := helper.HandleErrorResponse(resp, http.StatusAccepted, "create server")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var createdServer models.ServerBaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&createdServer); err != nil {
		fmt.Printf("JSON Decode Error: %v\n", err)
		return nil, err
	}

	return &createdServer, nil
}

func (s *ApiServerService) UpdateServer(id string, request *models.ServerUpdateRequest) (*models.ServerBaseResponse, error) {
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

	errorResponse := helper.HandleErrorResponse(resp, http.StatusOK, "update server")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var updatedServer models.ServerBaseResponse
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

	errorResponse := helper.HandleErrorResponse(resp, http.StatusAccepted, "delete server")
	if errorResponse != nil {
		return errorResponse
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

	model, diags := models.NewServerResourceModelFromAPI(context.Background(), server)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting to model: %v", diags)
	}

	return model, nil
}
