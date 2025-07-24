package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"terraform-provider-arsys-baremetal/internal/util"
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
						"must be a valid ID (e.g., 4EFAD5836CE43ACA502FD5B99BEE44EF)",
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
