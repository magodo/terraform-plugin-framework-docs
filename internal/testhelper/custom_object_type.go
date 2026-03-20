package testhelper

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Ensure the implementation satisfies the expected interfaces
var _ basetypes.ObjectTypable = CustomObjectType{}
var _ attr.TypeWithMarkdownDescription = CustomObjectType{}

type CustomObjectType struct {
	basetypes.ObjectType
	// ... potentially other fields ...
}

func (t CustomObjectType) MarkdownDescription(context.Context) string {
	return "A custom object type."
}

func (t CustomObjectType) Equal(o attr.Type) bool {
	other, ok := o.(CustomObjectType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t CustomObjectType) Object() string {
	return "CustomObjectType"
}

func (t CustomObjectType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	// CustomObjectValue defined in the value type section
	value := CustomObjectValue{
		ObjectValue: in,
	}

	return value, nil
}

func (t CustomObjectType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.ObjectType.ValueFromTerraform(ctx, in)

	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.ObjectValue)

	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromObject(ctx, stringValue)

	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting ObjectValue to ObjectValuable: %v", diags)
	}

	return stringValuable, nil
}

func (t CustomObjectType) ValueType(ctx context.Context) attr.Value {
	// CustomObjectValue defined in the value type section
	return CustomObjectValue{}
}
