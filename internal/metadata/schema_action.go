package metadata

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type ActionSchema struct {
	Description string
	Deprecation string

	Fields Fields

	// Including nested attribute object or block object.
	Nested NestedFields
}

func NewActionSchema(ctx context.Context, sch schema.Schema) (schema ActionSchema, diags diag.Diagnostics) {
	fields := Fields{}
	nested := NestedFields{}

	attrFields, attrNested, odiags := newActionAttrFields(ctx, nil, sch.Attributes)
	diags.Append(odiags...)
	if diags.HasError() {
		return
	}
	maps.Copy(fields, attrFields)
	maps.Copy(nested, attrNested)

	blockFields, blockNested, odiags := newActionBlockFields(ctx, nil, sch.Blocks)
	diags.Append(odiags...)
	if diags.HasError() {
		return
	}
	maps.Copy(fields, blockFields)
	maps.Copy(nested, blockNested)

	schema = ActionSchema{
		Description: *DescriptionOf(sch),
		Deprecation: sch.GetDeprecationMessage(),
		Fields:      fields,
		Nested:      nested,
	}
	return
}

func newActionAttrFields(ctx context.Context, parents []string, attrs map[string]schema.Attribute) (fields Fields, nested NestedFields, diags diag.Diagnostics) {
	fields = Fields{}
	nested = NestedFields{}

	for name, attr := range attrs {
		var (
			field Field

			objectNested NestedFields
			objectDiags  diag.Diagnostics
		)

		switch attr := attr.(type) {
		case schema.BoolAttribute:
			field = Field{
				parents:     parents,
				name:        name,
				dataType:    DTBool,
				required:    attr.IsRequired(),
				optional:    attr.IsOptional(),
				description: DescriptionOf(attr),
				deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Bool) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:   attr.IsWriteOnly(),
			}
		case schema.Float32Attribute:
			field = Field{
				parents:     parents,
				name:        name,
				dataType:    DTFloat32,
				required:    attr.IsRequired(),
				optional:    attr.IsOptional(),
				description: DescriptionOf(attr),
				deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Float32) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:   attr.IsWriteOnly(),
			}
		case schema.Float64Attribute:
			field = Field{
				parents:     parents,
				name:        name,
				dataType:    DTFloat64,
				required:    attr.IsRequired(),
				optional:    attr.IsOptional(),
				description: DescriptionOf(attr),
				deprecation: attr.DeprecationMessage,
				validators:  MapSlice(attr.Validators, func(v validator.Float64) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:   attr.IsWriteOnly(),
			}
		case schema.Int32Attribute:
			field = Field{
				parents:     parents,
				name:        name,
				dataType:    DTInt32,
				required:    attr.IsRequired(),
				optional:    attr.IsOptional(),
				description: DescriptionOf(attr),
				deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Int32) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:   attr.IsWriteOnly(),
			}
		case schema.Int64Attribute:
			field = Field{
				parents:     parents,
				name:        name,
				dataType:    DTInt64,
				required:    attr.IsRequired(),
				optional:    attr.IsOptional(),
				description: DescriptionOf(attr),
				deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Int64) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:   attr.IsWriteOnly(),
			}
		case schema.NumberAttribute:
			field = Field{
				parents:     parents,
				name:        name,
				dataType:    DTNumber,
				required:    attr.IsRequired(),
				optional:    attr.IsOptional(),
				description: DescriptionOf(attr),
				deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Number) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:   attr.IsWriteOnly(),
			}
		case schema.StringAttribute:
			field = Field{
				parents:     parents,
				name:        name,
				dataType:    DTString,
				required:    attr.IsRequired(),
				optional:    attr.IsOptional(),
				description: DescriptionOf(attr),
				deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.String) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:   attr.IsWriteOnly(),
			}
		case schema.ListAttribute:
			field = Field{
				parents:     parents,
				name:        name,
				dataType:    DTList,
				required:    attr.IsRequired(),
				optional:    attr.IsOptional(),
				description: DescriptionOf(attr),
				deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.List) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:   attr.IsWriteOnly(),
			}
		case schema.MapAttribute:
			field = Field{
				parents:     parents,
				name:        name,
				dataType:    DTMap,
				required:    attr.IsRequired(),
				optional:    attr.IsOptional(),
				description: DescriptionOf(attr),
				deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Map) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:   attr.IsWriteOnly(),
			}
		case schema.SetAttribute:
			field = Field{
				parents:     parents,
				name:        name,
				dataType:    DTSet,
				required:    attr.IsRequired(),
				optional:    attr.IsOptional(),
				description: DescriptionOf(attr),
				deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Set) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:   attr.IsWriteOnly(),
			}
		case schema.DynamicAttribute:
			field = Field{
				parents:     parents,
				name:        name,
				dataType:    DTDynamic,
				required:    attr.IsRequired(),
				optional:    attr.IsOptional(),
				description: DescriptionOf(attr),
				deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Dynamic) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:   attr.IsWriteOnly(),
			}

		case schema.ObjectAttribute:
			field = Field{
				parents:     parents,
				name:        name,
				dataType:    DTObjectAttr,
				required:    attr.IsRequired(),
				optional:    attr.IsOptional(),
				description: DescriptionOf(attr),
				deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:   attr.IsWriteOnly(),
			}
			// NOTE: We don't look into the AttributeTypes for an ObjectAttribute as it doesn't contain useful information.
		case schema.SingleNestedAttribute:
			field = Field{
				parents:     parents,
				name:        name,
				dataType:    DTSingleNestedAttr,
				required:    attr.IsRequired(),
				optional:    attr.IsOptional(),
				description: DescriptionOf(attr),
				deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:   attr.IsWriteOnly(),
			}
			objectNested, objectDiags = newActionNestedAttrObjectFields(ctx, slices.Concat(parents, []string{name}), attr.GetNestedObject().(schema.NestedAttributeObject))
		case schema.SetNestedAttribute:
			field = Field{
				parents:     parents,
				name:        name,
				dataType:    DTSetNestedAttr,
				required:    attr.IsRequired(),
				optional:    attr.IsOptional(),
				description: DescriptionOf(attr),
				deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Set) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:   attr.IsWriteOnly(),
			}
			objectNested, objectDiags = newActionNestedAttrObjectFields(ctx, slices.Concat(parents, []string{name}), attr.GetNestedObject().(schema.NestedAttributeObject))
		case schema.MapNestedAttribute:
			field = Field{
				parents:     parents,
				name:        name,
				dataType:    DTMapNestedAttr,
				required:    attr.IsRequired(),
				optional:    attr.IsOptional(),
				description: DescriptionOf(attr),
				deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Map) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:   attr.IsWriteOnly(),
			}
			objectNested, objectDiags = newActionNestedAttrObjectFields(ctx, slices.Concat(parents, []string{name}), attr.GetNestedObject().(schema.NestedAttributeObject))
		case schema.ListNestedAttribute:
			field = Field{
				parents:     parents,
				name:        name,
				dataType:    DTListNestedAttr,
				required:    attr.IsRequired(),
				optional:    attr.IsOptional(),
				description: DescriptionOf(attr),
				deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.List) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:   attr.IsWriteOnly(),
			}
			objectNested, objectDiags = newActionNestedAttrObjectFields(ctx, slices.Concat(parents, []string{name}), attr.GetNestedObject().(schema.NestedAttributeObject))
		default:
			diags.AddError("unknown schema type", fmt.Sprintf("%T", attr))
			return
		}

		fields[name] = field

		diags = append(diags, objectDiags...)
		if diags.HasError() {
			return
		}
		maps.Copy(nested, objectNested)
	}

	return
}

func newActionNestedAttrObjectFields(ctx context.Context, parents []string, obj schema.NestedAttributeObject) (nested NestedFields, diags diag.Diagnostics) {
	nested = NestedFields{}

	attrFields, attrNested, attrDiags := newActionAttrFields(ctx, parents, obj.Attributes)
	diags.Append(attrDiags...)
	if diags.HasError() {
		return
	}

	nested[strings.Join(parents, ".")] = NestedField{
		validators: MapSlice(obj.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
		fields:     attrFields,
	}
	maps.Copy(nested, attrNested)
	return
}

func newActionBlockFields(ctx context.Context, parents []string, blks map[string]schema.Block) (fields Fields, nested NestedFields, diags diag.Diagnostics) {
	fields = Fields{}
	nested = NestedFields{}

	for name, blk := range blks {
		var field Field

		switch blk := blk.(type) {
		case schema.SingleNestedBlock:
			field = Field{
				parents:     parents,
				name:        name,
				dataType:    DTSingleNestedBlock,
				optional:    true, // Always regard a block as optional.
				description: DescriptionOf(blk),
				deprecation: blk.GetDeprecationMessage(),
				validators:  MapSlice(blk.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.ListNestedBlock:
			field = Field{
				parents:     parents,
				name:        name,
				dataType:    DTListNestedBlock,
				optional:    true, // Always regard a block as optional.
				description: DescriptionOf(blk),
				deprecation: blk.GetDeprecationMessage(),
				validators:  MapSlice(blk.Validators, func(v validator.List) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.SetNestedBlock:
			field = Field{
				parents:     parents,
				name:        name,
				dataType:    DTSetNestedBlock,
				optional:    true, // Always regard a block as optional.
				description: DescriptionOf(blk),
				deprecation: blk.GetDeprecationMessage(),
				validators:  MapSlice(blk.Validators, func(v validator.Set) string { return DescriptionCtxOf(ctx, v) }),
			}
		}

		objectNested, odiags := newActionNestedBlkObjectFields(ctx, slices.Concat(parents, []string{name}), blk.GetNestedObject().(schema.NestedBlockObject))
		diags = append(diags, odiags...)
		if diags.HasError() {
			return
		}

		fields[name] = field
		maps.Copy(nested, objectNested)
	}

	return
}

func newActionNestedBlkObjectFields(ctx context.Context, parents []string, obj schema.NestedBlockObject) (nested NestedFields, diags diag.Diagnostics) {
	attrFields, attrNested, attrDiags := newActionAttrFields(ctx, parents, obj.Attributes)
	diags.Append(attrDiags...)
	if diags.HasError() {
		return
	}

	blkFields, blkNested, attrDiags := newActionBlockFields(ctx, parents, obj.Blocks)
	diags.Append(attrDiags...)
	if diags.HasError() {
		return
	}

	fields := Fields{}
	maps.Copy(fields, attrFields)
	maps.Copy(fields, blkFields)

	nested = NestedFields{}
	maps.Copy(nested, attrNested)
	maps.Copy(nested, blkNested)

	nested[strings.Join(parents, ".")] = NestedField{
		validators: MapSlice(obj.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
		fields:     fields,
	}
	return
}
