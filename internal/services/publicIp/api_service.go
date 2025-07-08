package publicIp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"terraform-provider-arsys-baremetal/internal/client"
	"terraform-provider-arsys-baremetal/internal/models"
)

var _ ApiPublicIpServiceInterface = (*ApiPublicIpService)(nil)

type ApiPublicIpService struct {
	client *client.APIClient
}

type ApiPublicIpServiceInterface interface {
	GetPublicIp(id string) (*models.PublicIpResponse, error)
	GetPublicIps() ([]models.PublicIpResponse, error)
	CreatePublicIp(request *models.PublicIpCreateRequest) (*models.PublicIpResponse, error)
	UpdatePublicIp(id string, request *models.PublicIpUpdateRequest) (*models.PublicIpResponse, error)
	DeletePublicIp(id string) error
}

func NewApiPublicIpService(client *client.APIClient) *ApiPublicIpService {
	return &ApiPublicIpService{client: client}
}

func GetPublicIpService(m interface{}) ApiPublicIpServiceInterface {
	if service, ok := m.(ApiPublicIpServiceInterface); ok {
		return service
	}

	if apiClient, ok := m.(*client.APIClient); ok {
		return NewApiPublicIpService(apiClient)
	}

	return nil
}

func (s *ApiPublicIpService) GetPublicIp(id string) (*models.PublicIpResponse, error) {
	resp, err := s.client.Get(fmt.Sprintf("/public_ips/%s", id))

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
		return nil, fmt.Errorf("public ip not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var publicIp models.PublicIpResponse
	if err := json.NewDecoder(resp.Body).Decode(&publicIp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	return &publicIp, nil
}

func (s *ApiPublicIpService) GetPublicIps() ([]models.PublicIpResponse, error) {
	resp, err := s.client.Get("/public_ips")

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

	var publicIps []models.PublicIpResponse
	if err := json.NewDecoder(resp.Body).Decode(&publicIps); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return publicIps, nil
}

func (s *ApiPublicIpService) CreatePublicIp(request *models.PublicIpCreateRequest) (*models.PublicIpResponse, error) {
	resp, err := s.client.Post("/public_ips", request)
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
		return nil, fmt.Errorf("error creating public ip: %s", string(body))
	}

	var createdPublicIp models.PublicIpResponse
	if err := json.NewDecoder(resp.Body).Decode(&createdPublicIp); err != nil {
		fmt.Printf("JSON Decode Error: %v\n", err)
		return nil, err
	}

	return &createdPublicIp, nil
}

func (s *ApiPublicIpService) UpdatePublicIp(id string, request *models.PublicIpUpdateRequest) (*models.PublicIpResponse, error) {
	resp, err := s.client.Put(fmt.Sprintf("/public_ips/%s", id), &request)
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
		return nil, fmt.Errorf("error updating public ip: %s", string(body))
	}

	var updatedPublicIp models.PublicIpResponse
	if err := json.NewDecoder(resp.Body).Decode(&updatedPublicIp); err != nil {
		return nil, err
	}

	return &updatedPublicIp, nil
}

func (s *ApiPublicIpService) DeletePublicIp(id string) error {
	resp, err := s.client.Delete(fmt.Sprintf("/public_ips/%s", id))
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
		return fmt.Errorf("error deleting public ip: %s", string(body))
	}

	return nil
}
