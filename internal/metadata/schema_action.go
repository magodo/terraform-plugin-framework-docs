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

	Fields ActionFields

	// Including nested attribute object or block object.
	Nested ActionNestedFields
}

func NewActionSchema(ctx context.Context, sch schema.Schema) (schema ActionSchema, diags diag.Diagnostics) {
	fields := ActionFields{}
	nested := ActionNestedFields{}

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
		Description: DescriptionOf(sch),
		Deprecation: sch.GetDeprecationMessage(),
		Fields:      fields,
		Nested:      nested,
	}
	return
}

func newActionAttrFields(ctx context.Context, parents []string, attrs map[string]schema.Attribute) (fields ActionFields, nested ActionNestedFields, diags diag.Diagnostics) {
	fields = ActionFields{}
	nested = ActionNestedFields{}

	for name, attr := range attrs {
		var (
			field ActionField

			objectNested ActionNestedFields
			objectDiags  diag.Diagnostics
		)

		switch attr := attr.(type) {
		case schema.BoolAttribute:
			field = ActionField{
				Parents:     parents,
				Name:        name,
				DataType:    DTBool,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Bool) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:   attr.IsWriteOnly(),
			}
		case schema.Float32Attribute:
			field = ActionField{
				Parents:     parents,
				Name:        name,
				DataType:    DTFloat32,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Float32) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:   attr.IsWriteOnly(),
			}
		case schema.Float64Attribute:
			field = ActionField{
				Parents:     parents,
				Name:        name,
				DataType:    DTFloat64,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Description: DescriptionOf(attr),
				Deprecation: attr.DeprecationMessage,
				validators:  MapSlice(attr.Validators, func(v validator.Float64) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:   attr.IsWriteOnly(),
			}
		case schema.Int32Attribute:
			field = ActionField{
				Parents:     parents,
				Name:        name,
				DataType:    DTInt32,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Int32) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:   attr.IsWriteOnly(),
			}
		case schema.Int64Attribute:
			field = ActionField{
				Parents:     parents,
				Name:        name,
				DataType:    DTInt64,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Int64) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:   attr.IsWriteOnly(),
			}
		case schema.NumberAttribute:
			field = ActionField{
				Parents:     parents,
				Name:        name,
				DataType:    DTNumber,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Number) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:   attr.IsWriteOnly(),
			}
		case schema.StringAttribute:
			field = ActionField{
				Parents:     parents,
				Name:        name,
				DataType:    DTString,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.String) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:   attr.IsWriteOnly(),
			}
		case schema.ListAttribute:
			field = ActionField{
				Parents:     parents,
				Name:        name,
				DataType:    DTList,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.List) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:   attr.IsWriteOnly(),
			}
		case schema.MapAttribute:
			field = ActionField{
				Parents:     parents,
				Name:        name,
				DataType:    DTMap,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Map) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:   attr.IsWriteOnly(),
			}
		case schema.SetAttribute:
			field = ActionField{
				Parents:     parents,
				Name:        name,
				DataType:    DTSet,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Set) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:   attr.IsWriteOnly(),
			}
		case schema.DynamicAttribute:
			field = ActionField{
				Parents:     parents,
				Name:        name,
				DataType:    DTDynamic,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Dynamic) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:   attr.IsWriteOnly(),
			}

		case schema.ObjectAttribute:
			field = ActionField{
				Parents:     parents,
				Name:        name,
				DataType:    DTObjectAttr,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:   attr.IsWriteOnly(),
			}
			// NOTE: We don't look into the AttributeTypes for an ObjectAttribute as it doesn't contain useful information.
		case schema.SingleNestedAttribute:
			field = ActionField{
				Parents:     parents,
				Name:        name,
				DataType:    DTSingleNestedAttr,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:   attr.IsWriteOnly(),
			}
			objectNested, objectDiags = newActionNestedAttrObjectFields(ctx, slices.Concat(parents, []string{name}), attr.GetNestedObject().(schema.NestedAttributeObject))
		case schema.SetNestedAttribute:
			field = ActionField{
				Parents:     parents,
				Name:        name,
				DataType:    DTSetNestedAttr,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Set) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:   attr.IsWriteOnly(),
			}
			objectNested, objectDiags = newActionNestedAttrObjectFields(ctx, slices.Concat(parents, []string{name}), attr.GetNestedObject().(schema.NestedAttributeObject))
		case schema.MapNestedAttribute:
			field = ActionField{
				Parents:     parents,
				Name:        name,
				DataType:    DTMapNestedAttr,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Map) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:   attr.IsWriteOnly(),
			}
			objectNested, objectDiags = newActionNestedAttrObjectFields(ctx, slices.Concat(parents, []string{name}), attr.GetNestedObject().(schema.NestedAttributeObject))
		case schema.ListNestedAttribute:
			field = ActionField{
				Parents:     parents,
				Name:        name,
				DataType:    DTListNestedAttr,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.List) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:   attr.IsWriteOnly(),
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

func newActionNestedAttrObjectFields(ctx context.Context, parents []string, obj schema.NestedAttributeObject) (nested ActionNestedFields, diags diag.Diagnostics) {
	nested = ActionNestedFields{}

	attrFields, attrNested, attrDiags := newActionAttrFields(ctx, parents, obj.Attributes)
	diags.Append(attrDiags...)
	if diags.HasError() {
		return
	}

	nested[strings.Join(parents, ".")] = ActionNestedField{
		Validators: MapSlice(obj.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
		Fields:     attrFields,
	}
	maps.Copy(nested, attrNested)
	return
}

func newActionBlockFields(ctx context.Context, parents []string, blks map[string]schema.Block) (fields ActionFields, nested ActionNestedFields, diags diag.Diagnostics) {
	fields = ActionFields{}
	nested = ActionNestedFields{}

	for name, blk := range blks {
		var field ActionField

		switch blk := blk.(type) {
		case schema.SingleNestedBlock:
			field = ActionField{
				Parents:     parents,
				Name:        name,
				DataType:    DTSingleNestedBlock,
				Optional:    true, // Always regard a block as optional.
				Description: DescriptionOf(blk),
				Deprecation: blk.GetDeprecationMessage(),
				validators:  MapSlice(blk.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.ListNestedBlock:
			field = ActionField{
				Parents:     parents,
				Name:        name,
				DataType:    DTListNestedBlock,
				Optional:    true, // Always regard a block as optional.
				Description: DescriptionOf(blk),
				Deprecation: blk.GetDeprecationMessage(),
				validators:  MapSlice(blk.Validators, func(v validator.List) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.SetNestedBlock:
			field = ActionField{
				Parents:     parents,
				Name:        name,
				DataType:    DTSetNestedBlock,
				Optional:    true, // Always regard a block as optional.
				Description: DescriptionOf(blk),
				Deprecation: blk.GetDeprecationMessage(),
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

func newActionNestedBlkObjectFields(ctx context.Context, parents []string, obj schema.NestedBlockObject) (nested ActionNestedFields, diags diag.Diagnostics) {
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

	fields := ActionFields{}
	maps.Copy(fields, attrFields)
	maps.Copy(fields, blkFields)

	nested = ActionNestedFields{}
	maps.Copy(nested, attrNested)
	maps.Copy(nested, blkNested)

	nested[strings.Join(parents, ".")] = ActionNestedField{
		Validators: MapSlice(obj.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
		Fields:     fields,
	}
	return
}
