package privateNetwork

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

var _ ApiPrivateNetworkServiceInterface = (*ApiPrivateNetworkService)(nil)

type ApiPrivateNetworkService struct {
	client *client.APIClient
}

type ApiPrivateNetworkServiceInterface interface {
	GetPrivateNetwork(id string) (*models.PrivateNetworkResponse, error)
	GetPrivateNetworks() ([]models.PrivateNetworkResponse, error)
	CreatePrivateNetwork(request *models.PrivateNetworkCreateRequest) (*models.PrivateNetworkResponse, error)
	UpdatePrivateNetwork(id string, request *models.PrivateNetworkUpdateRequest) (*models.PrivateNetworkResponse, error)
	DeletePrivateNetwork(id string) error
	GetPrivateNetworkServer(id, serverId string) (*models.PrivateNetworkServerResponse, error)
	GetPrivateNetworkServers(id string) ([]models.PrivateNetworkServerResponse, error)
	CreatePrivateNetworkServers(id string, request *models.PrivateNetworkServerRequest) (*models.PrivateNetworkResponse, error)
	DeletePrivateNetworkServer(privateNetworkId, serverId string) (*models.PrivateNetworkResponse, error)
}

func NewApiPrivateNetworkService(client *client.APIClient) *ApiPrivateNetworkService {
	return &ApiPrivateNetworkService{client: client}
}

func GetPrivateNetworkService(m interface{}) ApiPrivateNetworkServiceInterface {
	if service, ok := m.(ApiPrivateNetworkServiceInterface); ok {
		return service
	}

	if apiClient, ok := m.(*client.APIClient); ok {
		return NewApiPrivateNetworkService(apiClient)
	}

	return nil
}

func (s *ApiPrivateNetworkService) GetPrivateNetwork(id string) (*models.PrivateNetworkResponse, error) {
	resp, err := s.client.Get(fmt.Sprintf("/private_networks/%s", id))

	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	errorResponse := util.HandleErrorResponse(resp, http.StatusOK, "get private network")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var privateNetwork models.PrivateNetworkResponse
	if err := json.NewDecoder(resp.Body).Decode(&privateNetwork); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	return &privateNetwork, nil
}

func (s *ApiPrivateNetworkService) GetPrivateNetworks() ([]models.PrivateNetworkResponse, error) {
	resp, err := s.client.Get("/private_networks")

	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	errorResponse := util.HandleErrorResponse(resp, http.StatusOK, "get private networks")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var privateNetworks []models.PrivateNetworkResponse
	if err := json.NewDecoder(resp.Body).Decode(&privateNetworks); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return privateNetworks, nil
}

func (s *ApiPrivateNetworkService) CreatePrivateNetwork(request *models.PrivateNetworkCreateRequest) (*models.PrivateNetworkResponse, error) {
	resp, err := s.client.Post("/private_networks", request)
	if err != nil {
		return nil, err

	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	errorResponse := util.HandleErrorResponse(resp, http.StatusCreated, "create private network")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var createdPrivateNetwork models.PrivateNetworkResponse
	if err := json.NewDecoder(resp.Body).Decode(&createdPrivateNetwork); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &createdPrivateNetwork, nil
}

func (s *ApiPrivateNetworkService) UpdatePrivateNetwork(id string, request *models.PrivateNetworkUpdateRequest) (*models.PrivateNetworkResponse, error) {
	resp, err := s.client.Put(fmt.Sprintf("/private_networks/%s", id), &request)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	errorResponse := util.HandleErrorResponse(resp, http.StatusOK, "update private network")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var updatedPrivateNetwork models.PrivateNetworkResponse
	if err := json.NewDecoder(resp.Body).Decode(&updatedPrivateNetwork); err != nil {
		return nil, err
	}

	return &updatedPrivateNetwork, nil
}

func (s *ApiPrivateNetworkService) DeletePrivateNetwork(id string) error {
	resp, err := s.client.Delete(fmt.Sprintf("/private_networks/%s", id))
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	errorResponse := util.HandleErrorResponse(resp, http.StatusOK, "delete private network")
	if errorResponse != nil {
		return errorResponse
	}

	return nil
}

func (s *ApiPrivateNetworkService) GetPrivateNetworkServer(privateNetworkId, serverId string) (*models.PrivateNetworkServerResponse, error) {
	resp, err := s.client.Get(fmt.Sprintf("/private_networks/%s/servers/%s", privateNetworkId, serverId))

	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	errorResponse := util.HandleErrorResponse(resp, http.StatusOK, "get private network server")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var privateNetworkServer models.PrivateNetworkServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&privateNetworkServer); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	return &privateNetworkServer, nil
}

func (s *ApiPrivateNetworkService) GetPrivateNetworkServers(id string) ([]models.PrivateNetworkServerResponse, error) {
	resp, err := s.client.Get(fmt.Sprintf("/private_networks/%s/servers", id))

	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	errorResponse := util.HandleErrorResponse(resp, http.StatusOK, "get private network servers")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var privateNetworkServers []models.PrivateNetworkServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&privateNetworkServers); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return privateNetworkServers, nil
}

func (s *ApiPrivateNetworkService) CreatePrivateNetworkServers(id string, request *models.PrivateNetworkServerRequest) (*models.PrivateNetworkResponse, error) {
	resp, err := s.client.Post(fmt.Sprintf("/private_networks/%s/servers", id), request)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	errorResponse := util.HandleErrorResponse(resp, http.StatusAccepted, "assign servers to private network")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var response models.PrivateNetworkResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

func (s *ApiPrivateNetworkService) DeletePrivateNetworkServer(privateNetworkId, serverId string) (*models.PrivateNetworkResponse, error) {
	resp, err := s.client.Delete(fmt.Sprintf("/private_networks/%s/servers/%s", privateNetworkId, serverId))
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	errorResponse := util.HandleErrorResponse(resp, http.StatusAccepted, "delete private network server")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var response models.PrivateNetworkResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

func (s *ApiPrivateNetworkService) GetResource(id string) (util.ResourceModel, error) {
	privateNetwork, err := s.GetPrivateNetwork(id)
	if err != nil {
		return nil, err
	}

	if privateNetwork == nil {
		return nil, nil
	}

	model, diags := models.NewPrivateNetworkModel(context.Background(), privateNetwork)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting to model: %v", diags)
	}

	return model, nil
}
