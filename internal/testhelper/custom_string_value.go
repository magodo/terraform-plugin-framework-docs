package testhelper

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces
var _ basetypes.StringValuable = CustomStringValue{}

type CustomStringValue struct {
	basetypes.StringValue
	// ... potentially other fields ...
}

func (v CustomStringValue) Equal(o attr.Value) bool {
	other, ok := o.(CustomStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v CustomStringValue) Type(ctx context.Context) attr.Type {
	// CustomStringType defined in the schema type section
	return CustomStringType{}
}
