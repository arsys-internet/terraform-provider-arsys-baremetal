package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-arsys-baremetal/internal/models/server"
)

type ServerListModel struct {
	ServerBaseModel
	Status     types.Object `tfsdk:"status"`
	Hypervisor types.String `tfsdk:"hypervisor"`
}

type ServersModel struct {
	ID      types.String `tfsdk:"id"`
	Servers types.List   `tfsdk:"servers"`
}

type ServerListResponse struct {
	ServerBaseResponse
	Hypervisor string                    `json:"hypervisor"`
	Status     server.StatusBaseResponse `json:"status"`
}

func serverListModelObjectType() types.ObjectType {
	baseAttrs := serverBaseModelObjectType().AttrTypes

	attrs := make(map[string]attr.Type)
	for k, v := range baseAttrs {
		attrs[k] = v
	}
	attrs["status"] = server.StatusBaseObjectType()
	attrs["hypervisor"] = types.StringType

	return types.ObjectType{AttrTypes: attrs}
}

func NewServers(ctx context.Context, serversResponse []ServerListResponse) (*ServersModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	model := &ServersModel{}
	model.ID = types.StringValue("servers")

	serverModels, listDiags := NewServerFromList(ctx, serversResponse)
	diags.Append(listDiags...)

	if !listDiags.HasError() {
		serversList, convertDiags := types.ListValueFrom(ctx, serverListModelObjectType(), serverModels)
		diags.Append(convertDiags...)
		if !convertDiags.HasError() {
			model.Servers = serversList
		}
	}

	return model, diags
}

func NewServerFromList(ctx context.Context, serversResponse []ServerListResponse) ([]ServerListModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if len(serversResponse) == 0 {
		return []ServerListModel{}, diags
	}

	serverModels := make([]ServerListModel, 0, len(serversResponse))

	for _, serverResp := range serversResponse {
		baseModel, baseDiags := newServerBaseModelFromResponse(ctx, &serverResp.ServerBaseResponse)
		diags.Append(baseDiags...)

		if !baseDiags.HasError() {
			listModel := ServerListModel{
				ServerBaseModel: *baseModel,
			}

			statusObj, statusDiags := server.NewStatusBaseObject(serverResp.Status)
			diags.Append(statusDiags...)
			if !statusDiags.HasError() {
				listModel.Status = statusObj
			}

			serverModels = append(serverModels, listModel)
		}
	}

	return serverModels, diags
}

func serverNestedAttributeObject() schema.NestedAttributeObject {
	existingSchema := ServerDataSourceSchema(context.Background())

	attributes := make(map[string]schema.Attribute)
	for name, attribute := range existingSchema.Attributes {
		if name == "id" {
			attributes[name] = schema.StringAttribute{
				Computed:    true,
				Description: "Server identifier",
			}
		} else if name == "status" {
			statusAttrs := server.StatusDetailDataSourceSchema()
			delete(statusAttrs, "percent")

			attributes[name] = schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Server status",
				Attributes:  statusAttrs,
			}
		} else {
			attributes[name] = attribute
		}
	}

	delete(attributes, "recovery_mode")
	delete(attributes, "recovery_image_os")
	delete(attributes, "recovery_user")
	delete(attributes, "recovery_password")

	return schema.NestedAttributeObject{
		Attributes: attributes,
	}
}

func ServersDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Data source for listing servers",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for this data source",
			},
			"servers": schema.ListNestedAttribute{
				Computed:     true,
				Description:  "List of servers",
				NestedObject: serverNestedAttributeObject(),
			},
		},
	}
}
