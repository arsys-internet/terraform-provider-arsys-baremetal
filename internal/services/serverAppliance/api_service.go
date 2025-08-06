package serverAppliance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"terraform-provider-arsys-baremetal/internal/client"
	"terraform-provider-arsys-baremetal/internal/models"
	"terraform-provider-arsys-baremetal/internal/util"
)

var _ ApiServerApplianceServiceInterface = (*ApiServerApplianceService)(nil)

type ApiServerApplianceService struct {
	client *client.APIClient
}

type ApiServerApplianceServiceInterface interface {
	GetServerAppliance(id string) (*models.ServerApplianceResponse, error)
	GetServerAppliances() ([]models.ServerApplianceResponse, error)
}

func NewApiServerApplianceService(client *client.APIClient) *ApiServerApplianceService {
	return &ApiServerApplianceService{client: client}
}

func GetServerApplianceService(m interface{}) ApiServerApplianceServiceInterface {
	if service, ok := m.(ApiServerApplianceServiceInterface); ok {
		return service
	}

	if apiClient, ok := m.(*client.APIClient); ok {
		return NewApiServerApplianceService(apiClient)
	}

	return nil
}

func (s *ApiServerApplianceService) GetServerAppliance(id string) (*models.ServerApplianceResponse, error) {
	resp, err := s.client.Get(fmt.Sprintf("/server_appliances/%s", id))

	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	errorResponse := util.HandleErrorResponse(resp, http.StatusOK, "get server appliance")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var serverAppliance models.ServerApplianceResponse
	if err := json.NewDecoder(resp.Body).Decode(&serverAppliance); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	return &serverAppliance, nil
}

func (s *ApiServerApplianceService) GetServerAppliances() ([]models.ServerApplianceResponse, error) {
	resp, err := s.client.Get("/server_appliances")

	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	errorResponse := util.HandleErrorResponse(resp, http.StatusOK, "get server appliances")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var serverAppliances []models.ServerApplianceResponse
	if err := json.NewDecoder(resp.Body).Decode(&serverAppliances); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return serverAppliances, nil
}
