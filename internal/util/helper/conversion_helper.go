package helper

import "github.com/hashicorp/terraform-plugin-framework/types"

func AssignStringPtr(target **string, source types.String) {
	if !source.IsNull() && source.ValueString() != "" {
		value := source.ValueString()
		*target = &value
	}
}

func AssignFloatPtr(target **float64, source types.Float64) {
	if !source.IsNull() && source.ValueFloat64() != 0 {
		value := source.ValueFloat64()
		*target = &value
	}
}

func AssignBoolPtr(target **bool, source types.Bool) {
	if !source.IsNull() && source.ValueBool() {
		value := source.ValueBool()
		*target = &value
	}
}

func AssignStringDirect(target *string, source types.String) {
	if !source.IsNull() && source.ValueString() != "" {
		*target = source.ValueString()
	}
}

func GetStringValue(attr interface{}) string {
	if attr == nil {
		return ""
	}
	str := attr.(types.String)
	if str.IsNull() {
		return ""
	}
	return str.ValueString()
}

func GetStringPtr(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
