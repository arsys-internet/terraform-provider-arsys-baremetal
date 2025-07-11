package publicNetwork

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

var _ ApiPublicNetworkServiceInterface = (*ApiPublicNetworkService)(nil)

type ApiPublicNetworkService struct {
	client *client.APIClient
}

type ApiPublicNetworkServiceInterface interface {
	GetPublicNetwork(id string) (*models.PublicNetworkResponse, error)
	GetPublicNetworks() ([]models.PublicNetworkResponse, error)
	CreatePublicNetwork(request *models.PublicNetworkCreateRequest) (*models.PublicNetworkResponse, error)
	UpdatePublicNetwork(id string, request *models.PublicNetworkUpdateRequest) (*models.PublicNetworkResponse, error)
	DeletePublicNetwork(id string) error
}

func NewApiPublicNetworkService(client *client.APIClient) *ApiPublicNetworkService {
	return &ApiPublicNetworkService{client: client}
}

func GetPublicNetworkService(m interface{}) ApiPublicNetworkServiceInterface {
	if service, ok := m.(ApiPublicNetworkServiceInterface); ok {
		return service
	}

	if apiClient, ok := m.(*client.APIClient); ok {
		return NewApiPublicNetworkService(apiClient)
	}

	return nil
}

func (s *ApiPublicNetworkService) GetPublicNetwork(id string) (*models.PublicNetworkResponse, error) {
	resp, err := s.client.Get(fmt.Sprintf("/public_networks/%s", id))

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
		return nil, fmt.Errorf("public network not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var publicIp models.PublicNetworkResponse
	if err := json.NewDecoder(resp.Body).Decode(&publicIp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	return &publicIp, nil
}

func (s *ApiPublicNetworkService) GetPublicNetworks() ([]models.PublicNetworkResponse, error) {
	resp, err := s.client.Get("/public_networks")

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

	var publicIps []models.PublicNetworkResponse
	if err := json.NewDecoder(resp.Body).Decode(&publicIps); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return publicIps, nil
}

func (s *ApiPublicNetworkService) CreatePublicNetwork(request *models.PublicNetworkCreateRequest) (*models.PublicNetworkResponse, error) {
	resp, err := s.client.Post("/public_networks", request)
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
		return nil, fmt.Errorf("error creating public network: %s", string(body))
	}

	var createdPublicNetwork models.PublicNetworkCreateResponse
	if err := json.NewDecoder(resp.Body).Decode(&createdPublicNetwork); err != nil {
		return nil, fmt.Errorf("JSON Decode Error: %w", err)
	}

	return &createdPublicNetwork.Data, nil
}

func (s *ApiPublicNetworkService) UpdatePublicNetwork(id string, request *models.PublicNetworkUpdateRequest) (*models.PublicNetworkResponse, error) {
	resp, err := s.client.Put(fmt.Sprintf("/public_networks/%s", id), request)
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
		return nil, fmt.Errorf("error updating public network: %s", string(body))
	}

	var updatedPublicNetwork models.PublicNetworkCreateResponse
	if err := json.NewDecoder(resp.Body).Decode(&updatedPublicNetwork); err != nil {
		return nil, err
	}

	return &updatedPublicNetwork.Data, nil
}

func (s *ApiPublicNetworkService) DeletePublicNetwork(id string) error {
	resp, err := s.client.Delete(fmt.Sprintf("/public_networks/%s", id))
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error deleting public network: %s", string(body))
	}

	return nil
}

func (s *ApiPublicNetworkService) GetResource(id string) (util.ResourceModel, error) {
	network, err := s.GetPublicNetwork(id)
	if err != nil {
		return nil, err
	}

	if network == nil {
		return nil, nil
	}

	model, diags := models.NewPublicNetworkModel(context.Background(), network)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting to model: %v", diags)
	}

	return model, nil
}
