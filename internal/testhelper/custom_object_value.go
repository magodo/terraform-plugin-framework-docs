package testhelper

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces
var _ basetypes.ObjectValuable = CustomObjectValue{}

type CustomObjectValue struct {
	basetypes.ObjectValue
	// ... potentially other fields ...
}

func (v CustomObjectValue) Equal(o attr.Value) bool {
	other, ok := o.(CustomObjectValue)

	if !ok {
		return false
	}

	return v.ObjectValue.Equal(other.ObjectValue)
}

func (v CustomObjectValue) Type(ctx context.Context) attr.Type {
	// CustomObjectType defined in the schema type section
	return CustomObjectType{}
}
