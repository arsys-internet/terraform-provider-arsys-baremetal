package models

import (
	"context"
	"regexp"
	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PublicNetworkServerRequest struct {
	Servers []string `json:"servers"`
}

type PublicNetworkServerResourceModel struct {
	PublicNetworkId types.String   `tfsdk:"public_network_id"`
	Servers         []types.String `tfsdk:"servers"`
	Id              types.String   `tfsdk:"id"`
}

func PublicNetworkServerSchema(_ context.Context) rschema.Schema {
	return rschema.Schema{
		Description: "Public network resource",
		Attributes: map[string]rschema.Attribute{
			"public_network_id": rschema.StringAttribute{
				Required:    true,
				Description: "Public network identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid public network ID",
					),
				},
			},
			"servers": rschema.ListAttribute{
				Required:    true,
				Description: "List of servers identifiers in the public network",
				ElementType: types.StringType,
			},
			"id": rschema.StringAttribute{
				Computed:    true,
				Description: "Internal ID for the resource",
			},
		},
	}
}
