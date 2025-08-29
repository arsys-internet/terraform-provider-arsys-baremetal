package sshKey

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

var _ ApiSshKeyServiceInterface = (*ApiSshKeyService)(nil)

type ApiSshKeyService struct {
	client *client.APIClient
}

type ApiSshKeyServiceInterface interface {
	GetSshKey(id string) (*models.SshKeyResponse, error)
	GetSshKeys() ([]models.SshKeyResponse, error)
	CreateSshKey(request *models.SshKeyCreateRequest) (*models.SshKeyResponse, error)
	UpdateSshKey(id string, request *models.SshKeyUpdateRequest) (*models.SshKeyResponse, error)
	DeleteSshKey(id string) error
}

func NewApiSshKeyService(client *client.APIClient) *ApiSshKeyService {
	return &ApiSshKeyService{client: client}
}

func GetSshKeyService(m interface{}) ApiSshKeyServiceInterface {
	if service, ok := m.(ApiSshKeyServiceInterface); ok {
		return service
	}

	if apiClient, ok := m.(*client.APIClient); ok {
		return NewApiSshKeyService(apiClient)
	}

	return nil
}

func (s *ApiSshKeyService) GetSshKey(id string) (*models.SshKeyResponse, error) {
	resp, err := s.client.Get(fmt.Sprintf("/ssh_keys/%s", id))

	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	errorResponse := util.HandleErrorResponse(resp, http.StatusOK, "get SSH key")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var sshKey models.SshKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&sshKey); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	return &sshKey, nil
}

func (s *ApiSshKeyService) GetSshKeys() ([]models.SshKeyResponse, error) {
	resp, err := s.client.Get("/ssh_keys")

	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	errorResponse := util.HandleErrorResponse(resp, http.StatusOK, "get public SSH keys")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var sshKeys []models.SshKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&sshKeys); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return sshKeys, nil
}

func (s *ApiSshKeyService) CreateSshKey(request *models.SshKeyCreateRequest) (*models.SshKeyResponse, error) {
	resp, err := s.client.Post("/ssh_keys", request)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	errorResponse := util.HandleErrorResponse(resp, http.StatusCreated, "create public SSH key")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var createdSshKey models.SshKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&createdSshKey); err != nil {
		return nil, fmt.Errorf("JSON Decode Error: %w", err)
	}

	return &createdSshKey, nil
}

func (s *ApiSshKeyService) UpdateSshKey(id string, request *models.SshKeyUpdateRequest) (*models.SshKeyResponse, error) {
	resp, err := s.client.Put(fmt.Sprintf("/ssh_keys/%s", id), request)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	errorResponse := util.HandleErrorResponse(resp, http.StatusOK, "update public SSH key")
	if errorResponse != nil {
		return nil, errorResponse
	}

	var updatedSshKey models.SshKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&updatedSshKey); err != nil {
		return nil, err
	}

	return &updatedSshKey, nil
}

func (s *ApiSshKeyService) DeleteSshKey(id string) error {
	resp, err := s.client.Delete(fmt.Sprintf("/ssh_keys/%s", id))
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	errorResponse := util.HandleErrorResponse(resp, http.StatusOK, "delete public SSH key")
	if errorResponse != nil {
		return errorResponse
	}

	return nil
}

func (s *ApiSshKeyService) GetResource(id string) (util.ResourceModel, error) {
	ssh, err := s.GetSshKey(id)
	if err != nil {
		return nil, err
	}

	if ssh == nil {
		return nil, nil
	}

	model, diags := models.NewSshKeyResourceModel(context.Background(), ssh)
	if diags.HasError() {
		return nil, fmt.Errorf("error converting to model: %v", diags)
	}

	return model, nil
}
