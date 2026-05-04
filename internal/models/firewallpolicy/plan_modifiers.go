package firewallpolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type rulesWriteOnceModifier struct{}

// RulesWriteOnce uses the config value at creation time and the state value
// on subsequent updates. This prevents drift when rules are managed externally
// via firewall_policy_rule_add / firewall_policy_rule_remove resources, and
// avoids "Provider produced inconsistent result" errors caused by the provider
// returning the real API rule count instead of the planned one.
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
	var id types.String
	req.State.GetAttribute(ctx, path.Root("id"), &id)

	// During creation the resource has no id yet: use the config/plan value.
	if id.IsNull() || id.IsUnknown() || id.ValueString() == "" {
		return
	}

	// During updates: keep the current state value so Terraform does not plan
	// rule changes through this resource, and the post-apply result matches
	// what the plan promised.
	resp.PlanValue = req.StateValue
}
