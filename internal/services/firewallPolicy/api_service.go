package firewallPolicy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"terraform-provider-arsys-baremetal/internal/client"
	"terraform-provider-arsys-baremetal/internal/models"
	"terraform-provider-arsys-baremetal/internal/util"
)

var _ ApiFirewallPolicyServiceInterface = (*ApiFirewallPolicyService)(nil)

type ApiFirewallPolicyService struct {
	client *client.APIClient
}

type ApiFirewallPolicyServiceInterface interface {
	GetFirewallPolicy(id string) (*models.FirewallPoliciesResponse, error)
	GetFirewallPolicies() ([]models.FirewallPoliciesResponse, error)
	//CreateServer(request *models.ServerCreateRequest) (*models.ServerBaseResponse, error)
	//UpdateServer(id string, request *models.ServerUpdateRequest) (*models.ServerBaseResponse, error)
	//DeleteServer(id string) error
}

func NewAApiFirewallPolicyService(client *client.APIClient) *ApiFirewallPolicyService {
	return &ApiFirewallPolicyService{client: client}
}

func GetFirewallPolicyService(m interface{}) ApiFirewallPolicyServiceInterface {
	if service, ok := m.(ApiFirewallPolicyServiceInterface); ok {
		return service
	}

	if apiClient, ok := m.(*client.APIClient); ok {
		return NewAApiFirewallPolicyService(apiClient)
	}

	return nil
}

func (s *ApiFirewallPolicyService) GetFirewallPolicy(id string) (*models.FirewallPoliciesResponse, error) {
	resp, err := s.client.Get(fmt.Sprintf("/firewall_policies/%s", id))

	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	errorResponse := util.HandleErrorResponse(resp, http.StatusOK, "get firewall policy")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var server models.FirewallPoliciesResponse
	if err := json.NewDecoder(resp.Body).Decode(&server); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &server, nil
}

func (s *ApiFirewallPolicyService) GetFirewallPolicies() ([]models.FirewallPoliciesResponse, error) {
	resp, err := s.client.Get("/firewall_policies")

	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	errorResponse := util.HandleErrorResponse(resp, http.StatusOK, "get firewall policies")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var firewallPolicies []models.FirewallPoliciesResponse
	if err := json.NewDecoder(resp.Body).Decode(&firewallPolicies); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return firewallPolicies, nil
}

//func (s *ApiFirewallPolicyService) CreateServer(request *models.ServerCreateRequest) (*models.ServerBaseResponse, error) {
//	resp, err := s.client.Post("/servers", request)
//	if err != nil {
//		return nil, err
//	}
//
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			fmt.Println(err)
//		}
//	}(resp.Body)
//
//	errorResponse := util.HandleErrorResponse(resp, http.StatusAccepted, "create server")
//	if errorResponse != nil {
//		return nil, errorResponse
//	}
//
//	var createdServer models.ServerBaseResponse
//	if err := json.NewDecoder(resp.Body).Decode(&createdServer); err != nil {
//		fmt.Printf("JSON Decode Error: %v\n", err)
//		return nil, err
//	}
//
//	return &createdServer, nil
//}
//
//func (s *ApiFirewallPolicyService) UpdateServer(id string, request *models.ServerUpdateRequest) (*models.ServerBaseResponse, error) {
//	resp, err := s.client.Put(fmt.Sprintf("/servers/%s", id), request)
//	if err != nil {
//		return nil, err
//	}
//
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			fmt.Println(err)
//		}
//	}(resp.Body)
//
//	errorResponse := util.HandleErrorResponse(resp, http.StatusOK, "update server")
//	if errorResponse != nil {
//		return nil, errorResponse
//	}
//
//	var updatedServer models.ServerBaseResponse
//	if err := json.NewDecoder(resp.Body).Decode(&updatedServer); err != nil {
//		fmt.Printf("JSON Decode Error: %v\n", err)
//		return nil, err
//	}
//
//	return &updatedServer, nil
//}
//
//func (s *ApiFirewallPolicyService) DeleteServer(id string) error {
//	resp, err := s.client.Delete(fmt.Sprintf("/servers/%s", id))
//	if err != nil {
//		return err
//	}
//
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			fmt.Println(err)
//		}
//	}(resp.Body)
//
//	errorResponse := util.HandleErrorResponse(resp, http.StatusAccepted, "delete server")
//	if errorResponse != nil {
//		return errorResponse
//	}
//
//	return nil
//}
//

//func (s *ApiFirewallPolicyService) GetResource(id string) (util.ResourceModel, error) {
//	server, err := s.GetFirewallPolicy(id)
//	if err != nil {
//		return nil, err
//	}
//
//	if server == nil {
//		return nil, nil
//	}
//
//	model, diags := models.NewServerResourceModelFromAPI(context.Background(), server)
//	if diags.HasError() {
//		return nil, fmt.Errorf("error converting to model: %v", diags)
//	}
//
//	return model, nil
//}
