package datacenter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"terraform-provider-arsys-baremetal/internal/client"
	"terraform-provider-arsys-baremetal/internal/models"
	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ ApiDatacenterServiceInterface = (*ApiDatacenterService)(nil)

type ApiDatacenterService struct {
	client *client.APIClient
}

type ApiDatacenterServiceInterface interface {
	GetDatacenter(id string) (*models.DatacenterResponse, error)
	GetDatacenters() ([]models.DatacentersResponse, error)
}

func NewApiDatacenterService(client *client.APIClient) *ApiDatacenterService {
	return &ApiDatacenterService{client: client}
}

func GetDatacenterService(m interface{}) ApiDatacenterServiceInterface {
	if service, ok := m.(ApiDatacenterServiceInterface); ok {
		return service
	}

	if apiClient, ok := m.(*client.APIClient); ok {
		return NewApiDatacenterService(apiClient)
	}

	return nil
}

func (s *ApiDatacenterService) GetDatacenter(id string) (*models.DatacenterResponse, error) {
	resp, err := s.client.Get(fmt.Sprintf("/datacenters/%s", id))

	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			tflog.Warn(context.Background(), "Failed to close response body", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}(resp.Body)

	errorResponse := util.HandleErrorResponse(resp, http.StatusOK, "get datacenter")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var datacenter models.DatacenterResponse
	if err := json.NewDecoder(resp.Body).Decode(&datacenter); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	return &datacenter, nil
}

func (s *ApiDatacenterService) GetDatacenters() ([]models.DatacentersResponse, error) {
	resp, err := s.client.Get("/datacenters")

	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			tflog.Warn(context.Background(), "Failed to close response body", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}(resp.Body)

	errorResponse := util.HandleErrorResponse(resp, http.StatusOK, "get datacenters")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var datacenters []models.DatacentersResponse
	if err := json.NewDecoder(resp.Body).Decode(&datacenters); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return datacenters, nil
}
