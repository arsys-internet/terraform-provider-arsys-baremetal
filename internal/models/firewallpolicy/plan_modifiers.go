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

	var id types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if id.IsNull() || id.IsUnknown() || id.ValueString() == "" {
		return
	}

	resp.PlanValue = req.StateValue
}
