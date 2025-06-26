package privateNetwork

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"terraform-provider-arsys-baremetal/internal/client"
	"terraform-provider-arsys-baremetal/internal/models"
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
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("private network not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
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

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error creating private network: %s", string(body))
	}

	var createdPrivateNetwork models.PrivateNetworkResponse
	if err := json.NewDecoder(resp.Body).Decode(&createdPrivateNetwork); err != nil {
		fmt.Printf("JSON Decode Error: %v\n", err)
		return nil, err
	}

	return &createdPrivateNetwork, nil
}

func (s *ApiPrivateNetworkService) UpdatePrivateNetwork(id string, request *models.PrivateNetworkUpdateRequest) (*models.PrivateNetworkResponse, error) {
	fmt.Printf("DEBUG UpdatePrivateNetwork: ID=%s, Request=%+v\n", id, request)

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

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error updating private network: %s", string(body))
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

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error deleting private network: %s", string(body))
	}

	return nil
}
