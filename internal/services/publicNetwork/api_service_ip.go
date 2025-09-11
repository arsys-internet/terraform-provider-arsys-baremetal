package publicNetwork

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"terraform-provider-arsys-baremetal/internal/client"
	"terraform-provider-arsys-baremetal/internal/models"
	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const CompositeIDSeparator = ":"

func CreateCompositeID(publicNetworkId, ipId string) string {
	return fmt.Sprintf("%s%s%s", publicNetworkId, CompositeIDSeparator, ipId)
}

func ParseCompositeID(compositeId string) (publicNetworkId, ipId string, err error) {
	parts := strings.Split(compositeId, CompositeIDSeparator)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid composite ID format: %s, expected format: publicNetworkId:ipId", compositeId)
	}
	return parts[0], parts[1], nil
}

var _ ApiPublicNetworkIpServiceInterface = (*ApiPublicNetworkIpService)(nil)

type ApiPublicNetworkIpService struct {
	client *client.APIClient
}

type ApiPublicNetworkIpServiceInterface interface {
	GetPublicNetworkIp(publicNetworkId string, id string) (*models.PublicNetworkIpResponse, error)
	GetPublicNetworkIps(publicNetworkId string) ([]models.PublicNetworkIpResponse, error)
	AssignIpToPublicNetwork(id string, request *models.PublicNetworkIpRequest) (*[]models.PublicNetworkIpResponse, error)
	GetResource(id string) (util.ResourceModel, error)
}

func NewApiPublicNetworkIpService(client *client.APIClient) *ApiPublicNetworkIpService {
	return &ApiPublicNetworkIpService{client: client}
}

func GetPublicNetworkIpService(m interface{}) ApiPublicNetworkIpServiceInterface {
	if service, ok := m.(ApiPublicNetworkIpServiceInterface); ok {
		return service
	}

	if apiClient, ok := m.(*client.APIClient); ok {
		return NewApiPublicNetworkIpService(apiClient)
	}

	return nil
}

func (s *ApiPublicNetworkIpService) GetPublicNetworkIp(publicNetworkId string, id string) (*models.PublicNetworkIpResponse, error) {
	resp, err := s.client.Get(fmt.Sprintf("/public_networks/%s/ips/%s", publicNetworkId, id))

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

	errorResponse := util.HandleErrorResponse(resp, http.StatusOK, "get ips in the public network")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var publicIp models.PublicNetworkIpResponse
	if err := json.NewDecoder(resp.Body).Decode(&publicIp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	return &publicIp, nil
}

func (s *ApiPublicNetworkIpService) GetPublicNetworkIps(publicNetworkId string) ([]models.PublicNetworkIpResponse, error) {
	resp, err := s.client.Get(fmt.Sprintf("/public_networks/%s/ips", publicNetworkId))

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

	errorResponse := util.HandleErrorResponse(resp, http.StatusOK, "get ips in the public network")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var publicIps []models.PublicNetworkIpResponse
	if err := json.NewDecoder(resp.Body).Decode(&publicIps); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return publicIps, nil
}

func (s *ApiPublicNetworkIpService) AssignIpToPublicNetwork(id string, request *models.PublicNetworkIpRequest) (*[]models.PublicNetworkIpResponse, error) {
	resp, err := s.client.Put(fmt.Sprintf("/public_networks/%s/access", id), request)
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

	errorResponse := util.HandleErrorResponse(resp, http.StatusAccepted, "assign ips to public network")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var response models.PublicNetworkIpCreateResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	return &response.Data.Items, nil

}

func (s *ApiPublicNetworkIpService) GetResource(compositeId string) (util.ResourceModel, error) {
	publicNetworkId, ipId, err := ParseCompositeID(compositeId)
	if err != nil {
		return nil, err
	}

	ip, err := s.GetPublicNetworkIp(publicNetworkId, ipId)
	if err != nil {
		return nil, err
	}

	if ip == nil {
		return nil, nil
	}

	model, diags := models.NewPublicNetworkIpModel(context.Background(), ip)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting to model: %v", diags)
	}

	return model, nil
}
