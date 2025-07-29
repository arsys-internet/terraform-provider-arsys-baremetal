package models

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FirewallPoliciesModel struct {
	ID               types.String `tfsdk:"id"`
	FirewallPolicies types.List   `tfsdk:"firewall_policies"`
}

type FirewallPoliciesResponse = FirewallPolicyResponse

func NewFirewallPoliciesFromList(ctx context.Context, fpList []FirewallPolicyResponse) ([]FirewallPolicyModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	var models []FirewallPolicyModel

	if len(fpList) == 0 {
		return []FirewallPolicyModel{}, diags
	}

	for i, fp := range fpList {
		model, modelDiags := NewFirewallPolicyModel(ctx, fp)
		if modelDiags.HasError() {
			diags.AddError(
				"Build error",
				fmt.Sprintf("Failed to create firewall policy model for item %d: %s", i, modelDiags.Errors()[0].Summary()),
			)
			continue
		}
		diags.Append(modelDiags...)
		models = append(models, *model)
	}

	return models, diags
}

func NewFirewallPolicies(ctx context.Context, fpResponse []FirewallPolicyResponse) (*FirewallPoliciesModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	model := &FirewallPoliciesModel{}
	model.ID = types.StringValue("firewall_policies")

	firewallPolicies, listDiags := NewFirewallPoliciesFromList(ctx, fpResponse)
	diags.Append(listDiags...)

	if !listDiags.HasError() {
		firewallPoliciesList, convertDiags := types.ListValueFrom(ctx, FirewallPolicyObjectType(), firewallPolicies)
		diags.Append(convertDiags...)
		if !convertDiags.HasError() {
			model.FirewallPolicies = firewallPoliciesList
		}
	}

	return model, diags
}

func firewallPolicyNestedAttributeObject() schema.NestedAttributeObject {
	existingSchema := FirewallPolicyDataSourceSchema(context.Background())

	attributes := make(map[string]schema.Attribute)
	for name, attribute := range existingSchema.Attributes {
		if name == "id" {
			attributes[name] = schema.StringAttribute{
				Computed:    true,
				Description: "Firewall policy identifier",
			}
		} else {
			attributes[name] = attribute
		}
	}

	return schema.NestedAttributeObject{
		Attributes: attributes,
	}
}

func FirewallPoliciesDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Data source for listing firewall policies",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for this data source",
			},
			"firewall_policies": schema.ListNestedAttribute{
				Computed:     true,
				Description:  "List of firewall policies",
				NestedObject: firewallPolicyNestedAttributeObject(),
			},
		},
	}
}
