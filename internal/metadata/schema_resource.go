package metadata

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type ResourceSchema struct {
	Description string
	Deprecation string

	Fields Fields

	// Including nested attribute object or block object.
	Nested NestedFields
}

func NewResourceSchema(ctx context.Context, sch schema.Schema) (schema ResourceSchema, diags diag.Diagnostics) {
	fields := Fields{}
	nested := NestedFields{}

	attrFields, attrNested, odiags := newResourceAttrFields(ctx, nil, sch.Attributes)
	diags.Append(odiags...)
	if diags.HasError() {
		return
	}
	maps.Copy(fields, attrFields)
	maps.Copy(nested, attrNested)

	blockFields, blockNested, odiags := newResourceBlockFields(ctx, nil, sch.Blocks)
	diags.Append(odiags...)
	if diags.HasError() {
		return
	}
	maps.Copy(fields, blockFields)
	maps.Copy(nested, blockNested)

	schema = ResourceSchema{
		Description: DescriptionOf(sch),
		Deprecation: sch.GetDeprecationMessage(),
		Fields:      fields,
		Nested:      nested,
	}
	return
}

func newResourceAttrFields(ctx context.Context, parents []string, attrs map[string]schema.Attribute) (fields Fields, nested NestedFields, diags diag.Diagnostics) {
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
				parents:       parents,
				name:          name,
				dataType:      DTBool,
				required:      attr.IsRequired(),
				optional:      attr.IsOptional(),
				computed:      attr.IsComputed(),
				sensitive:     attr.IsSensitive(),
				description:   DescriptionOf(attr),
				deprecation:   attr.GetDeprecationMessage(),
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Bool) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.Bool) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrZero(attr.Default, func(v defaults.Bool) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:     attr.IsWriteOnly(),
			}
		case schema.Float32Attribute:
			field = Field{
				parents:       parents,
				name:          name,
				dataType:      DTFloat32,
				required:      attr.IsRequired(),
				optional:      attr.IsOptional(),
				computed:      attr.IsComputed(),
				sensitive:     attr.IsSensitive(),
				description:   DescriptionOf(attr),
				deprecation:   attr.GetDeprecationMessage(),
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Float32) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.Float32) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrZero(attr.Default, func(v defaults.Float32) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:     attr.IsWriteOnly(),
			}
		case schema.Float64Attribute:
			field = Field{
				parents:       parents,
				name:          name,
				dataType:      DTFloat64,
				required:      attr.IsRequired(),
				optional:      attr.IsOptional(),
				computed:      attr.IsComputed(),
				sensitive:     attr.IsSensitive(),
				description:   DescriptionOf(attr),
				deprecation:   attr.DeprecationMessage,
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Float64) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.Float64) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrZero(attr.Default, func(v defaults.Float64) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:     attr.IsWriteOnly(),
			}
		case schema.Int32Attribute:
			field = Field{
				parents:       parents,
				name:          name,
				dataType:      DTInt32,
				required:      attr.IsRequired(),
				optional:      attr.IsOptional(),
				computed:      attr.IsComputed(),
				sensitive:     attr.IsSensitive(),
				description:   DescriptionOf(attr),
				deprecation:   attr.GetDeprecationMessage(),
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Int32) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.Int32) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrZero(attr.Default, func(v defaults.Int32) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:     attr.IsWriteOnly(),
			}
		case schema.Int64Attribute:
			field = Field{
				parents:       parents,
				name:          name,
				dataType:      DTInt64,
				required:      attr.IsRequired(),
				optional:      attr.IsOptional(),
				computed:      attr.IsComputed(),
				sensitive:     attr.IsSensitive(),
				description:   DescriptionOf(attr),
				deprecation:   attr.GetDeprecationMessage(),
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Int64) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.Int64) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrZero(attr.Default, func(v defaults.Int64) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:     attr.IsWriteOnly(),
			}
		case schema.NumberAttribute:
			field = Field{
				parents:       parents,
				name:          name,
				dataType:      DTNumber,
				required:      attr.IsRequired(),
				optional:      attr.IsOptional(),
				computed:      attr.IsComputed(),
				sensitive:     attr.IsSensitive(),
				description:   DescriptionOf(attr),
				deprecation:   attr.GetDeprecationMessage(),
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Number) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.Number) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrZero(attr.Default, func(v defaults.Number) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:     attr.IsWriteOnly(),
			}
		case schema.StringAttribute:
			field = Field{
				parents:       parents,
				name:          name,
				dataType:      DTString,
				required:      attr.IsRequired(),
				optional:      attr.IsOptional(),
				computed:      attr.IsComputed(),
				sensitive:     attr.IsSensitive(),
				description:   DescriptionOf(attr),
				deprecation:   attr.GetDeprecationMessage(),
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.String) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.String) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrZero(attr.Default, func(v defaults.String) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:     attr.IsWriteOnly(),
			}
		case schema.ListAttribute:
			field = Field{
				parents:       parents,
				name:          name,
				dataType:      DTList,
				required:      attr.IsRequired(),
				optional:      attr.IsOptional(),
				computed:      attr.IsComputed(),
				sensitive:     attr.IsSensitive(),
				description:   DescriptionOf(attr),
				deprecation:   attr.GetDeprecationMessage(),
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.List) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.List) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrZero(attr.Default, func(v defaults.List) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:     attr.IsWriteOnly(),
			}
		case schema.MapAttribute:
			field = Field{
				parents:       parents,
				name:          name,
				dataType:      DTMap,
				required:      attr.IsRequired(),
				optional:      attr.IsOptional(),
				computed:      attr.IsComputed(),
				sensitive:     attr.IsSensitive(),
				description:   DescriptionOf(attr),
				deprecation:   attr.GetDeprecationMessage(),
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Map) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.Map) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrZero(attr.Default, func(v defaults.Map) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:     attr.IsWriteOnly(),
			}
		case schema.SetAttribute:
			field = Field{
				parents:       parents,
				name:          name,
				dataType:      DTSet,
				required:      attr.IsRequired(),
				optional:      attr.IsOptional(),
				computed:      attr.IsComputed(),
				sensitive:     attr.IsSensitive(),
				description:   DescriptionOf(attr),
				deprecation:   attr.GetDeprecationMessage(),
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Set) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.Set) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrZero(attr.Default, func(v defaults.Set) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:     attr.IsWriteOnly(),
			}
		case schema.DynamicAttribute:
			field = Field{
				parents:       parents,
				name:          name,
				dataType:      DTDynamic,
				required:      attr.IsRequired(),
				optional:      attr.IsOptional(),
				computed:      attr.IsComputed(),
				sensitive:     attr.IsSensitive(),
				description:   DescriptionOf(attr),
				deprecation:   attr.GetDeprecationMessage(),
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Dynamic) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.Dynamic) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrZero(attr.Default, func(v defaults.Dynamic) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:     attr.IsWriteOnly(),
			}

		case schema.ObjectAttribute:
			field = Field{
				parents:       parents,
				name:          name,
				dataType:      DTObjectAttr,
				required:      attr.IsRequired(),
				optional:      attr.IsOptional(),
				computed:      attr.IsComputed(),
				sensitive:     attr.IsSensitive(),
				description:   DescriptionOf(attr),
				deprecation:   attr.GetDeprecationMessage(),
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Object) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrZero(attr.Default, func(v defaults.Object) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:     attr.IsWriteOnly(),
			}
			objects, objectDiags := newObjects(ctx, slices.Concat(parents, []string{name}), attr.AttributeTypes)
			if objectDiags.HasError() {
				return nil, nil, objectDiags
			}
			objectNested = objects.ToNestedFields(field)
		case schema.SingleNestedAttribute:
			field = Field{
				parents:       parents,
				name:          name,
				dataType:      DTSingleNestedAttr,
				required:      attr.IsRequired(),
				optional:      attr.IsOptional(),
				computed:      attr.IsComputed(),
				sensitive:     attr.IsSensitive(),
				description:   DescriptionOf(attr),
				deprecation:   attr.GetDeprecationMessage(),
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Object) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrZero(attr.Default, func(v defaults.Object) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:     attr.IsWriteOnly(),
			}
			objectNested, objectDiags = newResourceNestedAttrObjectFields(ctx, slices.Concat(parents, []string{name}), attr.GetNestedObject().(schema.NestedAttributeObject))
		case schema.SetNestedAttribute:
			field = Field{
				parents:       parents,
				name:          name,
				dataType:      DTSetNestedAttr,
				required:      attr.IsRequired(),
				optional:      attr.IsOptional(),
				computed:      attr.IsComputed(),
				sensitive:     attr.IsSensitive(),
				description:   DescriptionOf(attr),
				deprecation:   attr.GetDeprecationMessage(),
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Set) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.Set) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrZero(attr.Default, func(v defaults.Set) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:     attr.IsWriteOnly(),
			}
			objectNested, objectDiags = newResourceNestedAttrObjectFields(ctx, slices.Concat(parents, []string{name}), attr.GetNestedObject().(schema.NestedAttributeObject))
		case schema.MapNestedAttribute:
			field = Field{
				parents:       parents,
				name:          name,
				dataType:      DTMapNestedAttr,
				required:      attr.IsRequired(),
				optional:      attr.IsOptional(),
				computed:      attr.IsComputed(),
				sensitive:     attr.IsSensitive(),
				description:   DescriptionOf(attr),
				deprecation:   attr.GetDeprecationMessage(),
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Map) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.Map) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrZero(attr.Default, func(v defaults.Map) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:     attr.IsWriteOnly(),
			}
			objectNested, objectDiags = newResourceNestedAttrObjectFields(ctx, slices.Concat(parents, []string{name}), attr.GetNestedObject().(schema.NestedAttributeObject))
		case schema.ListNestedAttribute:
			field = Field{
				parents:       parents,
				name:          name,
				dataType:      DTListNestedAttr,
				required:      attr.IsRequired(),
				optional:      attr.IsOptional(),
				computed:      attr.IsComputed(),
				sensitive:     attr.IsSensitive(),
				description:   DescriptionOf(attr),
				deprecation:   attr.GetDeprecationMessage(),
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.List) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.List) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrZero(attr.Default, func(v defaults.List) string { return DescriptionCtxOf(ctx, v) }),
				writeOnly:     attr.IsWriteOnly(),
			}
			objectNested, objectDiags = newResourceNestedAttrObjectFields(ctx, slices.Concat(parents, []string{name}), attr.GetNestedObject().(schema.NestedAttributeObject))
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

func newResourceNestedAttrObjectFields(ctx context.Context, parents []string, obj schema.NestedAttributeObject) (nested NestedFields, diags diag.Diagnostics) {
	nested = NestedFields{}

	attrFields, attrNested, attrDiags := newResourceAttrFields(ctx, parents, obj.Attributes)
	diags.Append(attrDiags...)
	if diags.HasError() {
		return
	}

	nested[strings.Join(parents, ".")] = NestedField{
		planModifiers: MapSlice(obj.PlanModifiers, func(v planmodifier.Object) string { return DescriptionCtxOf(ctx, v) }),
		validators:    MapSlice(obj.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
		fields:        attrFields,
	}
	maps.Copy(nested, attrNested)
	return
}

func newResourceBlockFields(ctx context.Context, parents []string, blks map[string]schema.Block) (fields Fields, nested NestedFields, diags diag.Diagnostics) {
	fields = Fields{}
	nested = NestedFields{}

	for name, blk := range blks {
		var field Field

		switch blk := blk.(type) {
		case schema.SingleNestedBlock:
			field = Field{
				parents:       parents,
				name:          name,
				dataType:      DTSingleNestedBlock,
				optional:      true, // Always regard a block as optional.
				description:   DescriptionOf(blk),
				deprecation:   blk.GetDeprecationMessage(),
				planModifiers: MapSlice(blk.PlanModifiers, func(v planmodifier.Object) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(blk.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.ListNestedBlock:
			field = Field{
				parents:       parents,
				name:          name,
				dataType:      DTListNestedBlock,
				optional:      true, // Always regard a block as optional.
				description:   DescriptionOf(blk),
				deprecation:   blk.GetDeprecationMessage(),
				planModifiers: MapSlice(blk.PlanModifiers, func(v planmodifier.List) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(blk.Validators, func(v validator.List) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.SetNestedBlock:
			field = Field{
				parents:       parents,
				name:          name,
				dataType:      DTSetNestedBlock,
				optional:      true, // Always regard a block as optional.
				description:   DescriptionOf(blk),
				deprecation:   blk.GetDeprecationMessage(),
				planModifiers: MapSlice(blk.PlanModifiers, func(v planmodifier.Set) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(blk.Validators, func(v validator.Set) string { return DescriptionCtxOf(ctx, v) }),
			}
		}

		objectNested, odiags := newResourceNestedBlkObjectFields(ctx, slices.Concat(parents, []string{name}), blk.GetNestedObject().(schema.NestedBlockObject))
		diags = append(diags, odiags...)
		if diags.HasError() {
			return
		}

		fields[name] = field
		maps.Copy(nested, objectNested)
	}

	return
}

func newResourceNestedBlkObjectFields(ctx context.Context, parents []string, obj schema.NestedBlockObject) (nested NestedFields, diags diag.Diagnostics) {
	attrFields, attrNested, attrDiags := newResourceAttrFields(ctx, parents, obj.Attributes)
	diags.Append(attrDiags...)
	if diags.HasError() {
		return
	}

	blkFields, blkNested, attrDiags := newResourceBlockFields(ctx, parents, obj.Blocks)
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
		planModifiers: MapSlice(obj.PlanModifiers, func(v planmodifier.Object) string { return DescriptionCtxOf(ctx, v) }),
		validators:    MapSlice(obj.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
		fields:        fields,
	}
	return
}
