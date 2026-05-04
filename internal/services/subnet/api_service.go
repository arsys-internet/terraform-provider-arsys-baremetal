package subnet

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

var _ ApiSubnetServiceInterface = (*ApiSubnetService)(nil)

type ApiSubnetService struct {
	client *client.APIClient
}

type ApiSubnetServiceInterface interface {
	GetSubnet(id string) (*models.SubnetResponse, error)
	CreateSubnet(request *models.CreateSubnetRequest) (*models.SubnetResponse, error)
	DeleteSubnet(id string) error
}

func NewApiSubnetService(client *client.APIClient) *ApiSubnetService {
	return &ApiSubnetService{client: client}
}

func GetSubnetService(m interface{}) ApiSubnetServiceInterface {
	if service, ok := m.(ApiSubnetServiceInterface); ok {
		return service
	}

	if apiClient, ok := m.(*client.APIClient); ok {
		return NewApiSubnetService(apiClient)
	}

	return nil
}

func (s *ApiSubnetService) GetSubnet(id string) (*models.SubnetResponse, error) {
	resp, err := s.client.Get(fmt.Sprintf("/subnets/%s", id))
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

	errorResponse := util.HandleErrorResponse(resp, http.StatusOK, "get subnet")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var subnet models.SubnetResponse
	if err := json.NewDecoder(resp.Body).Decode(&subnet); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &subnet, nil
}

func (s *ApiSubnetService) CreateSubnet(request *models.CreateSubnetRequest) (*models.SubnetResponse, error) {
	resp, err := s.client.Post("/subnets", request)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			tflog.Warn(context.Background(), "Failed to close response body", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}(resp.Body)

	errorResponse := util.HandleErrorResponse(resp, http.StatusCreated, "create subnet")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var createdSubnet models.SubnetResponse
	if err := json.NewDecoder(resp.Body).Decode(&createdSubnet); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &createdSubnet, nil
}

func (s *ApiSubnetService) DeleteSubnet(id string) error {
	resp, err := s.client.Delete(fmt.Sprintf("/subnets/%s", id))
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			tflog.Warn(context.Background(), "Failed to close response body", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}(resp.Body)

	errorResponse := util.HandleErrorResponse(resp, http.StatusOK, "delete subnet")
	if errorResponse != nil {
		return errorResponse
	}

	return nil
}
