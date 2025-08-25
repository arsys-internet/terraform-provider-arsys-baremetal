package publicIp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"terraform-provider-arsys-baremetal/internal/client"
	"terraform-provider-arsys-baremetal/internal/models"
	"terraform-provider-arsys-baremetal/internal/util"
)

var _ ApiPublicIpServiceInterface = (*ApiPublicIpService)(nil)

type ApiPublicIpService struct {
	client *client.APIClient
}

var subnetTypeRegex = regexp.MustCompile(`^IPV[46]`)

type ApiPublicIpServiceInterface interface {
	GetPublicIp(id string) (*models.PublicIpResponse, error)
	GetPublicIps() ([]models.PublicIpResponse, error)
	CreatePublicIp(request *models.PublicIpCreateRequest) (*models.PublicIpResponse, error)
	UpdatePublicIp(id string, request *models.PublicIpUpdateRequest) (*models.PublicIpResponse, error)
	DeletePublicIp(id string) error
	GetSubnets() ([]models.PublicIpResponse, error)
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

	errorResponse := util.HandleErrorResponse(resp, http.StatusOK, "get public ip")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var publicIp models.PublicIpResponse
	if err := json.NewDecoder(resp.Body).Decode(&publicIp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if subnetTypeRegex.MatchString(publicIp.Type) {
		return nil, fmt.Errorf("public ip not found")
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

	errorResponse := util.HandleErrorResponse(resp, http.StatusOK, "get public ips")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var publicIps []models.PublicIpResponse
	if err := json.NewDecoder(resp.Body).Decode(&publicIps); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	var filteredPublicIps []models.PublicIpResponse
	for _, ip := range publicIps {
		if subnetTypeRegex.MatchString(ip.Type) {
			filteredPublicIps = append(filteredPublicIps, ip)
		}
	}

	return filteredPublicIps, nil
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

	errorResponse := util.HandleErrorResponse(resp, http.StatusCreated, "create public ip")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var createdPublicIp models.PublicIpResponse
	if err := json.NewDecoder(resp.Body).Decode(&createdPublicIp); err != nil {
		return nil, fmt.Errorf("JSON Decode Error: %w", err)
	}

	return &createdPublicIp, nil
}

func (s *ApiPublicIpService) UpdatePublicIp(id string, request *models.PublicIpUpdateRequest) (*models.PublicIpResponse, error) {
	resp, err := s.client.Put(fmt.Sprintf("/public_ips/%s", id), request)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	errorResponse := util.HandleErrorResponse(resp, http.StatusOK, "update public ip")
	if errorResponse != nil {
		return nil, errorResponse
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

	errorResponse := util.HandleErrorResponse(resp, http.StatusAccepted, "delete public ip")
	if errorResponse != nil {
		return errorResponse
	}

	return nil
}

func (s *ApiPublicIpService) GetSubnets() ([]models.PublicIpResponse, error) {
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

	errorResponse := util.HandleErrorResponse(resp, http.StatusOK, "get subnets")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var publicIps []models.PublicIpResponse
	if err := json.NewDecoder(resp.Body).Decode(&publicIps); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	var filteredPublicIps []models.PublicIpResponse
	for _, ip := range publicIps {
		subnetRegex := regexp.MustCompile(`^IPV[46]Subnet$`)
		if subnetRegex.MatchString(ip.Type) {
			filteredPublicIps = append(filteredPublicIps, ip)
		}
	}

	return filteredPublicIps, nil
}

func (s *ApiPublicIpService) GetResource(id string) (util.ResourceModel, error) {
	ip, err := s.GetPublicIp(id)
	if err != nil {
		return nil, err
	}

	if ip == nil {
		return nil, nil
	}

	model, diags := models.NewPublicIpResourceModel(context.Background(), ip)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting to model: %v", diags)
	}

	return model, nil
}
