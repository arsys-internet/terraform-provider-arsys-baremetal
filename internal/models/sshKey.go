package models

import (
	"context"
	"fmt"
	"regexp"
	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SshKeyModel struct {
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	State        types.String `tfsdk:"state"`
	Servers      types.List   `tfsdk:"servers"`
	Md5          types.String `tfsdk:"md5"`
	PublicKey    types.String `tfsdk:"public_key"`
	CreationDate types.String `tfsdk:"creation_date"`
	PrivateKey   types.String `tfsdk:"private_key"`
}

type SshKeyResponse struct {
	Id           string               `json:"id"`
	Name         string               `json:"name"`
	Description  *string              `json:"description"`
	State        string               `json:"state"`
	Servers      []IdentifierResponse `json:"servers"`
	Md5          string               `json:"md5"`
	PublicKey    string               `json:"public_key"`
	CreationDate string               `json:"creation_date"`
	PrivateKey   *string              `json:"private_key,omitempty"`
}

func newSshKeyFromResponse(_ context.Context, ssh *SshKeyResponse) (*SshKeyModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if ssh == nil {
		diags.AddError("Response Error", "SSH key response is nil")
		return nil, diags
	}

	model := &SshKeyModel{}

	model.Id = types.StringValue(ssh.Id)
	model.Name = types.StringValue(ssh.Name)
	if ssh.Description != nil {
		model.Description = types.StringValue(*ssh.Description)
	} else {
		model.Description = types.StringNull()
	}
	model.State = types.StringValue(ssh.State)
	model.Md5 = types.StringValue(ssh.Md5)
	model.PublicKey = types.StringValue(ssh.PublicKey)
	model.CreationDate = types.StringValue(ssh.CreationDate)

	if ssh.PrivateKey != nil {
		model.PrivateKey = types.StringValue(*ssh.PrivateKey)
	} else {
		model.PrivateKey = types.StringNull()
	}

	serversList, listDiags := NewIdentifierList(ssh.Servers)
	diags.Append(listDiags...)
	if !listDiags.HasError() {
		model.Servers = serversList
	}

	return model, diags
}

func NewSshKeyModel(ctx context.Context, ssh *SshKeyResponse) (*SshKeyModel, diag.Diagnostics) {
	return newSshKeyFromResponse(ctx, ssh)
}

func NewSshKeyFromList(ctx context.Context, sshList []SshKeyResponse) ([]SshKeyModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	var models []SshKeyModel

	if len(sshList) == 0 {
		return []SshKeyModel{}, diags
	}

	for i, ssh := range sshList {
		model, modelDiags := NewSshKeyModel(ctx, &ssh)
		if modelDiags.HasError() {
			diags.AddError(
				"List Constructor Error",
				fmt.Sprintf("failed to create model for item %d: %s", i, modelDiags.Errors()[0].Summary()),
			)
			continue
		}
		diags.Append(modelDiags...)
		if model != nil {
			models = append(models, *model)
		}
	}

	return models, diags
}

func SshKeyDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Data source for SSH key information",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "SSH key identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid ID (e.g., 4EFAD5836CE43ACA502FD5B99BEE44EF)",
					),
				},
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "SSH key name",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "SSH key description",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "Current state of the SSH key",
			},
			"servers": schema.ListNestedAttribute{
				Computed:     true,
				Description:  "List of servers associated with the SSH key",
				NestedObject: IdentifierNestedObject(),
			},
			"md5": schema.StringAttribute{
				Computed:    true,
				Description: "MD5 hash of the SSH key",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid MD5 hash (32 hexadecimal characters)",
					),
				},
			},
			"public_key": schema.StringAttribute{
				Computed:    true,
				Description: "SSH public key content",
			},
			"creation_date": schema.StringAttribute{
				Computed:    true,
				Description: "SSH key creation date in ISO 8601 format",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.DateTimePattern),
						"must be a date in ISO 8601 format (e.g., 2023-05-29T09:43:31+00:00)",
					),
				},
			},
			"private_key": schema.StringAttribute{
				Computed:    true,
				Description: "SSH key private key",
			},
		},
	}
}

type SshKeyCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	PublicKey   string `json:"public_key,omitempty"`
}

func (m *SshKeyModel) ToCreateRequest() SshKeyCreateRequest {
	return SshKeyCreateRequest{
		Name:        m.Name.ValueString(),
		Description: m.Description.ValueString(),
		PublicKey:   m.PublicKey.ValueString(),
	}
}

func NewSshKeyResourceModel(ctx context.Context, ssh *SshKeyResponse) (*SshKeyModel, diag.Diagnostics) {
	baseModel, diags := newSshKeyFromResponse(ctx, ssh)
	if diags.HasError() {
		return nil, diags
	}

	return baseModel, diags
}

func SshKeyResourceSchema(_ context.Context) rschema.Schema {
	return rschema.Schema{
		Description: "SSH key resource",
		Attributes: map[string]rschema.Attribute{
			"id": rschema.StringAttribute{
				Computed:    true,
				Description: "SSH key identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": rschema.StringAttribute{
				Required:    true,
				Description: "SSH key name",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(util.MaxNameLength),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": rschema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "SSH key description",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(util.MaxDescriptionLength),
				},
			},
			"state": rschema.StringAttribute{
				Computed:    true,
				Description: "Current state of the SSH key",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"servers": rschema.ListNestedAttribute{
				Computed:     true,
				Description:  "List of servers associated with the SSH key",
				NestedObject: IdentifierResourceNestedObject(),
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"md5": rschema.StringAttribute{
				Computed:    true,
				Description: "MD5 hash of the SSH key",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid MD5 hash (32 hexadecimal characters)",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"public_key": rschema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "SSH public key content",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(util.MaxNameLength),
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"creation_date": rschema.StringAttribute{
				Computed:    true,
				Description: "SSH key creation date in ISO 8601 format",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.DateTimePattern),
						"must be a date in ISO 8601 format (e.g., 2023-05-29T09:43:31+00:00)",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"private_key": rschema.StringAttribute{
				Computed:    true,
				Description: "SSH key private key",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

type SshKeyUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

func (m *SshKeyModel) ToUpdateRequest() SshKeyUpdateRequest {
	return SshKeyUpdateRequest{
		Name:        m.Name.ValueString(),
		Description: m.Description.ValueString(),
	}
}

func (m *SshKeyModel) GetState() string {
	return m.State.ValueString()
}
