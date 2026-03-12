package metadata

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type ListSchema struct {
	Description string
	Deprecation string

	Fields ListFields

	// Including nested attribute object or block object.
	Nested ListNestedFields
}

func NewListSchema(ctx context.Context, sch schema.Schema) (schema ListSchema, diags diag.Diagnostics) {
	fields := ListFields{}
	nested := ListNestedFields{}

	attrFields, attrNested, odiags := newListAttrFields(ctx, nil, sch.Attributes)
	diags.Append(odiags...)
	if diags.HasError() {
		return
	}
	maps.Copy(fields, attrFields)
	maps.Copy(nested, attrNested)

	blockFields, blockNested, odiags := newListBlockFields(ctx, nil, sch.Blocks)
	diags.Append(odiags...)
	if diags.HasError() {
		return
	}
	maps.Copy(fields, blockFields)
	maps.Copy(nested, blockNested)

	schema = ListSchema{
		Description: DescriptionOf(sch),
		Deprecation: sch.GetDeprecationMessage(),
		Fields:      fields,
		Nested:      nested,
	}
	return
}

func newListAttrFields(ctx context.Context, parents []string, attrs map[string]schema.Attribute) (fields ListFields, nested ListNestedFields, diags diag.Diagnostics) {
	fields = ListFields{}
	nested = ListNestedFields{}

	for name, attr := range attrs {
		var (
			field ListField

			objectNested ListNestedFields
			objectDiags  diag.Diagnostics
		)

		switch attr := attr.(type) {
		case schema.BoolAttribute:
			field = ListField{
				Parents:     parents,
				Name:        name,
				DataType:    DTBool,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Bool) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.Float32Attribute:
			field = ListField{
				Parents:     parents,
				Name:        name,
				DataType:    DTFloat32,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Float32) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.Float64Attribute:
			field = ListField{
				Parents:     parents,
				Name:        name,
				DataType:    DTFloat64,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Description: DescriptionOf(attr),
				Deprecation: attr.DeprecationMessage,
				validators:  MapSlice(attr.Validators, func(v validator.Float64) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.Int32Attribute:
			field = ListField{
				Parents:     parents,
				Name:        name,
				DataType:    DTInt32,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Int32) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.Int64Attribute:
			field = ListField{
				Parents:     parents,
				Name:        name,
				DataType:    DTInt64,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Int64) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.NumberAttribute:
			field = ListField{
				Parents:     parents,
				Name:        name,
				DataType:    DTNumber,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Number) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.StringAttribute:
			field = ListField{
				Parents:     parents,
				Name:        name,
				DataType:    DTString,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.String) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.ListAttribute:
			field = ListField{
				Parents:     parents,
				Name:        name,
				DataType:    DTList,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.List) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.MapAttribute:
			field = ListField{
				Parents:     parents,
				Name:        name,
				DataType:    DTMap,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Map) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.DynamicAttribute:
			field = ListField{
				Parents:     parents,
				Name:        name,
				DataType:    DTDynamic,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Dynamic) string { return DescriptionCtxOf(ctx, v) }),
			}

		case schema.ObjectAttribute:
			field = ListField{
				Parents:     parents,
				Name:        name,
				DataType:    DTObjectAttr,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
			}
			// NOTE: We don't look into the AttributeTypes for an ObjectAttribute as it doesn't contain useful information.
		case schema.SingleNestedAttribute:
			field = ListField{
				Parents:     parents,
				Name:        name,
				DataType:    DTSingleNestedAttr,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
			}
			objectNested, objectDiags = newListNestedAttrObjectFields(ctx, slices.Concat(parents, []string{name}), attr.GetNestedObject().(schema.NestedAttributeObject))
		case schema.MapNestedAttribute:
			field = ListField{
				Parents:     parents,
				Name:        name,
				DataType:    DTMapNestedAttr,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Map) string { return DescriptionCtxOf(ctx, v) }),
			}
			objectNested, objectDiags = newListNestedAttrObjectFields(ctx, slices.Concat(parents, []string{name}), attr.GetNestedObject().(schema.NestedAttributeObject))
		case schema.ListNestedAttribute:
			field = ListField{
				Parents:     parents,
				Name:        name,
				DataType:    DTListNestedAttr,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.List) string { return DescriptionCtxOf(ctx, v) }),
			}
			objectNested, objectDiags = newListNestedAttrObjectFields(ctx, slices.Concat(parents, []string{name}), attr.GetNestedObject().(schema.NestedAttributeObject))
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

func newListNestedAttrObjectFields(ctx context.Context, parents []string, obj schema.NestedAttributeObject) (nested ListNestedFields, diags diag.Diagnostics) {
	nested = ListNestedFields{}

	attrFields, attrNested, attrDiags := newListAttrFields(ctx, parents, obj.Attributes)
	diags.Append(attrDiags...)
	if diags.HasError() {
		return
	}

	nested[strings.Join(parents, ".")] = ListNestedField{
		Validators: MapSlice(obj.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
		Fields:     attrFields,
	}
	maps.Copy(nested, attrNested)
	return
}

func newListBlockFields(ctx context.Context, parents []string, blks map[string]schema.Block) (fields ListFields, nested ListNestedFields, diags diag.Diagnostics) {
	fields = ListFields{}
	nested = ListNestedFields{}

	for name, blk := range blks {
		var field ListField

		switch blk := blk.(type) {
		case schema.SingleNestedBlock:
			field = ListField{
				Parents:     parents,
				Name:        name,
				DataType:    DTSingleNestedBlock,
				Optional:    true, // Always regard a block as optional.
				Description: DescriptionOf(blk),
				Deprecation: blk.GetDeprecationMessage(),
				validators:  MapSlice(blk.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.ListNestedBlock:
			field = ListField{
				Parents:     parents,
				Name:        name,
				DataType:    DTListNestedBlock,
				Optional:    true, // Always regard a block as optional.
				Description: DescriptionOf(blk),
				Deprecation: blk.GetDeprecationMessage(),
				validators:  MapSlice(blk.Validators, func(v validator.List) string { return DescriptionCtxOf(ctx, v) }),
			}
		}

		objectNested, odiags := newListNestedBlkObjectFields(ctx, slices.Concat(parents, []string{name}), blk.GetNestedObject().(schema.NestedBlockObject))
		diags = append(diags, odiags...)
		if diags.HasError() {
			return
		}

		fields[name] = field
		maps.Copy(nested, objectNested)
	}

	return
}

func newListNestedBlkObjectFields(ctx context.Context, parents []string, obj schema.NestedBlockObject) (nested ListNestedFields, diags diag.Diagnostics) {
	attrFields, attrNested, attrDiags := newListAttrFields(ctx, parents, obj.Attributes)
	diags.Append(attrDiags...)
	if diags.HasError() {
		return
	}

	blkFields, blkNested, attrDiags := newListBlockFields(ctx, parents, obj.Blocks)
	diags.Append(attrDiags...)
	if diags.HasError() {
		return
	}

	fields := ListFields{}
	maps.Copy(fields, attrFields)
	maps.Copy(fields, blkFields)

	nested = ListNestedFields{}
	maps.Copy(nested, attrNested)
	maps.Copy(nested, blkNested)

	nested[strings.Join(parents, ".")] = ListNestedField{
		Validators: MapSlice(obj.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
		Fields:     fields,
	}
	return
}
