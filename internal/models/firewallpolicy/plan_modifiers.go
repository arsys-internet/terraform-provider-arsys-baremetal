package firewallpolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type rulesWriteOnceModifier struct{}

func RulesWriteOnce() planmodifier.List {
	return rulesWriteOnceModifier{}
}

func (m rulesWriteOnceModifier) Description(_ context.Context) string {
	return "Rules are applied only at creation. After creation the provider reflects the API state."
}

func (m rulesWriteOnceModifier) MarkdownDescription(_ context.Context) string {
	return "Rules are applied only at creation. After creation the provider reflects the API state."
}

func (m rulesWriteOnceModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if req.State.Raw.IsNull() || !req.State.Raw.IsKnown() {
		return
	}

	var stateID types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &stateID)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if stateID.IsNull() || stateID.IsUnknown() || stateID.ValueString() == "" {
		return
	}

	var planID types.String
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("id"), &planID)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if planID.IsNull() || planID.IsUnknown() || planID.ValueString() == "" {
		return
	}

	resp.PlanValue = req.StateValue
}
