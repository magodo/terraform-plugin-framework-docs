package metadata

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
)

type ResourceIdentitySchema struct {
	Fields ResourceIdentityFields
}

func NewResourceIdentitySchema(ctx context.Context, sch identityschema.Schema) (schema ResourceIdentitySchema, diags diag.Diagnostics) {
	fields, odiags := newResourceIdentityAttrFields(ctx, sch.Attributes)
	diags.Append(odiags...)
	if diags.HasError() {
		return
	}
	schema = ResourceIdentitySchema{
		Fields: fields,
	}
	return
}

func newResourceIdentityAttrFields(_ context.Context, attrs map[string]identityschema.Attribute) (fields ResourceIdentityFields, diags diag.Diagnostics) {
	fields = ResourceIdentityFields{}

	for name, attr := range attrs {
		var field ResourceIdentityField

		switch attr := attr.(type) {
		case identityschema.BoolAttribute:
			field = ResourceIdentityField{
				Name:        name,
				DataType:    DTBool,
				Required:    attr.IsRequiredForImport(),
				Optional:    attr.IsOptionalForImport(),
				Description: DescriptionOf(attr),
			}
		case identityschema.Float32Attribute:
			field = ResourceIdentityField{
				Name:        name,
				DataType:    DTFloat32,
				Required:    attr.IsRequiredForImport(),
				Optional:    attr.IsOptionalForImport(),
				Description: DescriptionOf(attr),
			}
		case identityschema.Float64Attribute:
			field = ResourceIdentityField{
				Name:        name,
				DataType:    DTFloat64,
				Required:    attr.IsRequiredForImport(),
				Optional:    attr.IsOptionalForImport(),
				Description: DescriptionOf(attr),
			}
		case identityschema.Int32Attribute:
			field = ResourceIdentityField{
				Name:        name,
				DataType:    DTInt32,
				Required:    attr.IsRequiredForImport(),
				Optional:    attr.IsOptionalForImport(),
				Description: DescriptionOf(attr),
			}
		case identityschema.Int64Attribute:
			field = ResourceIdentityField{
				Name:        name,
				DataType:    DTInt64,
				Required:    attr.IsRequiredForImport(),
				Optional:    attr.IsOptionalForImport(),
				Description: DescriptionOf(attr),
			}
		case identityschema.NumberAttribute:
			field = ResourceIdentityField{
				Name:        name,
				DataType:    DTNumber,
				Required:    attr.IsRequiredForImport(),
				Optional:    attr.IsOptionalForImport(),
				Description: DescriptionOf(attr),
			}
		case identityschema.StringAttribute:
			field = ResourceIdentityField{
				Name:        name,
				DataType:    DTString,
				Required:    attr.IsRequiredForImport(),
				Optional:    attr.IsOptionalForImport(),
				Description: DescriptionOf(attr),
			}
		case identityschema.ListAttribute:
			field = ResourceIdentityField{
				Name:        name,
				DataType:    DTList,
				Required:    attr.IsRequiredForImport(),
				Optional:    attr.IsOptionalForImport(),
				Description: DescriptionOf(attr),
			}
		default:
			diags.AddError("unknown identity schema type", fmt.Sprintf("%T", attr))
			return
		}

		fields[name] = field
	}

	return
}
