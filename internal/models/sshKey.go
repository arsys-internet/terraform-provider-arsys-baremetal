package models

import (
	"context"
	"fmt"
	"regexp"
	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
				Description: "Current state of the IP (ACTIVE, etc.)",
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
		},
	}
}

//
//type SshKeyCreateRequest struct {
//	ReverseDns   string `json:"reverse_dns,omitempty"`
//	DatacenterId string `json:"datacenter_id,omitempty"`
//	Type         string `json:"type,omitempty"`
//}
//
//func (m *SshKeyResourceModel) ToCreateRequest() SshKeyCreateRequest {
//	return SshKeyCreateRequest{
//		ReverseDns:   m.ReverseDNS.ValueString(),
//		DatacenterId: m.DatacenterId.ValueString(),
//		Type:         m.Type.ValueString(),
//	}
//}
//
//type SshKeyResourceModel struct {
//	SshKeyModel
//	DatacenterId types.String `tfsdk:"datacenter_id"`
//}

//func NewSshKeyResourceModel(ctx context.Context, ssh *SshKeyResponse) (*SshKeyResourceModel, diag.Diagnostics) {
//	baseModel, diags := newSshKeyFromResponse(ctx, ssh)
//	if diags.HasError() {
//		return nil, diags
//	}
//
//	resourceModel := &SshKeyResourceModel{
//		SshKeyModel:  *baseModel,
//		DatacenterId: types.StringValue(ip.Datacenter.ID),
//	}
//
//	return resourceModel, diags
//}

//func SshKeyResourceSchema(_ context.Context) rschema.Schema {
//	return rschema.Schema{
//		Description: "Public ip resource",
//		Attributes: map[string]rschema.Attribute{
//			"id": rschema.StringAttribute{
//				Computed:    true,
//				Description: "Public ip identifier",
//			},
//			"reverse_dns": rschema.StringAttribute{
//				Optional:    true,
//				Computed:    true,
//				Description: "Reverse DNS name",
//				Validators: []validator.String{
//					stringvalidator.LengthAtMost(util.MaxNameLength),
//					stringvalidator.LengthAtLeast(1),
//				},
//			},
//			"datacenter_id": rschema.StringAttribute{
//				Optional:    true,
//				Computed:    true,
//				Description: "Datacenter identifier where the ip will be created",
//				Validators: []validator.String{
//					stringvalidator.RegexMatches(
//						regexp.MustCompile(util.HexID32Pattern),
//						"must be a valid datacenter_id (e.g., 4EEAD5836CF43ACA502FD5B99BFF44EF)",
//					),
//				},
//			},
//			"type": rschema.StringAttribute{
//				Optional:    true,
//				Computed:    true,
//				Description: "IP type",
//				Validators: []validator.String{
//					stringvalidator.OneOf("IPV4", "IPV6"),
//				},
//			},
//			"ip": rschema.StringAttribute{
//				Computed:    true,
//				Description: "IP address",
//			},
//			"datacenter":  BaseDatacenterNestedAttribute(),
//			"assigned_to": AssignedToNestedAttribute(),
//			"subnet_id": rschema.StringAttribute{
//				Computed:    true,
//				Description: "ID of the subnet to which the IP belongs",
//			},
//			"is_dhcp": rschema.BoolAttribute{
//				Computed:    true,
//				Description: "IP use DHCP",
//			},
//			"state": rschema.StringAttribute{
//				Computed:    true,
//				Description: "Current state of the IP (ACTIVE, etc.)",
//			},
//			"creation_date": rschema.StringAttribute{
//				Computed:    true,
//				Description: "IP creation date",
//			},
//		},
//	}
//}
//
//type SshKeyUpdateRequest struct {
//	ReverseDns string `json:"reverse_dns"`
//}
//
//func (m *SshKeyResourceModel) ToUpdateRequest() SshKeyUpdateRequest {
//	return SshKeyUpdateRequest{
//		ReverseDns: m.ReverseDNS.ValueString(),
//	}
//}

//func (m *SshKeyResourceModel) GetState() string {
//	return m.State.ValueString()
//}
