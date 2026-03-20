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
var _ basetypes.StringTypable = CustomStringType{}
var _ attr.TypeWithMarkdownDescription = CustomStringType{}

type CustomStringType struct {
	basetypes.StringType
	// ... potentially other fields ...
}

func (t CustomStringType) MarkdownDescription(context.Context) string {
	return "A custom string type."
}

func (t CustomStringType) Equal(o attr.Type) bool {
	other, ok := o.(CustomStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t CustomStringType) String() string {
	return "CustomStringType"
}

func (t CustomStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	// CustomStringValue defined in the value type section
	value := CustomStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t CustomStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)

	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)

	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)

	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}

func (t CustomStringType) ValueType(ctx context.Context) attr.Value {
	// CustomStringValue defined in the value type section
	return CustomStringValue{}
}
