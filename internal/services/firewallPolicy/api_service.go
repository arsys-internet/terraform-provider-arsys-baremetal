package firewallPolicy

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

var _ ApiFirewallPolicyServiceInterface = (*ApiFirewallPolicyService)(nil)

type ApiFirewallPolicyService struct {
	client *client.APIClient
}

type ApiFirewallPolicyServiceInterface interface {
	GetFirewallPolicy(id string) (*models.FirewallPolicyResponse, error)
	GetFirewallPolicies() ([]models.FirewallPolicyResponse, error)
	CreateFirewallPolicy(request *models.FirewallPolicyCreateRequest) (*models.FirewallPolicyResponse, error)
	UpdateFirewallPolicy(id string, request *models.FirewallPolicyUpdateRequest) (*models.FirewallPolicyResponse, error)
	DeleteFirewallPolicy(id string) error
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

func (s *ApiFirewallPolicyService) GetFirewallPolicy(id string) (*models.FirewallPolicyResponse, error) {
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

	var firewallPolicy models.FirewallPolicyResponse
	if err := json.NewDecoder(resp.Body).Decode(&firewallPolicy); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &firewallPolicy, nil
}

func (s *ApiFirewallPolicyService) GetFirewallPolicies() ([]models.FirewallPolicyResponse, error) {
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

	var firewallPolicies []models.FirewallPolicyResponse
	if err := json.NewDecoder(resp.Body).Decode(&firewallPolicies); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return firewallPolicies, nil
}

func (s *ApiFirewallPolicyService) CreateFirewallPolicy(request *models.FirewallPolicyCreateRequest) (*models.FirewallPolicyResponse, error) {
	resp, err := s.client.Post("/firewall_policies", request)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	errorResponse := util.HandleErrorResponse(resp, http.StatusAccepted, "create firewall policy")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var createdServer models.FirewallPolicyResponse
	if err := json.NewDecoder(resp.Body).Decode(&createdServer); err != nil {
		fmt.Printf("JSON Decode Error: %v\n", err)
		return nil, err
	}

	return &createdServer, nil
}

func (s *ApiFirewallPolicyService) UpdateFirewallPolicy(id string, request *models.FirewallPolicyUpdateRequest) (*models.FirewallPolicyResponse, error) {
	resp, err := s.client.Put(fmt.Sprintf("/firewall_policies/%s", id), request)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	errorResponse := util.HandleErrorResponse(resp, http.StatusOK, "update firewall policy")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var updatedServer models.FirewallPolicyResponse
	if err := json.NewDecoder(resp.Body).Decode(&updatedServer); err != nil {
		fmt.Printf("JSON Decode Error: %v\n", err)
		return nil, err
	}

	return &updatedServer, nil
}

func (s *ApiFirewallPolicyService) DeleteFirewallPolicy(id string) error {
	resp, err := s.client.Delete(fmt.Sprintf("/firewall_policies/%s", id))
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	errorResponse := util.HandleErrorResponse(resp, http.StatusAccepted, "delete firewall policy")
	if errorResponse != nil {
		return errorResponse
	}

	return nil
}

func (s *ApiFirewallPolicyService) GetResource(id string) (util.ResourceModel, error) {
	firewallPolicy, err := s.GetFirewallPolicy(id)
	if err != nil {
		return nil, err
	}

	if firewallPolicy == nil {
		return nil, nil
	}

	model, diags := models.NewFirewallPolicyModel(context.Background(), *firewallPolicy)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting to model: %v", diags)
	}

	return model, nil
}
