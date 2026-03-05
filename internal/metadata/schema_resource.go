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

	Fields ResourceFields

	// Including nested attribute object or block object.
	Nested ResourceNestedFields
}

func NewResourceSchema(ctx context.Context, sch schema.Schema) (schema ResourceSchema, diags diag.Diagnostics) {
	fields := ResourceFields{}
	nested := ResourceNestedFields{}

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

func newResourceAttrFields(ctx context.Context, parents []string, attrs map[string]schema.Attribute) (fields ResourceFields, nested ResourceNestedFields, diags diag.Diagnostics) {
	fields = ResourceFields{}
	nested = ResourceNestedFields{}

	for name, attr := range attrs {
		var (
			field ResourceField

			objectNested ResourceNestedFields
			objectDiags  diag.Diagnostics
		)

		switch attr := attr.(type) {
		case schema.BoolAttribute:
			field = ResourceField{
				Parents:       parents,
				Name:          name,
				DataType:      DTBool,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.GetDeprecationMessage(),
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Bool) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.Bool) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrNil(attr.Default, func(v defaults.Bool) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}
		case schema.Float32Attribute:
			field = ResourceField{
				Parents:       parents,
				Name:          name,
				DataType:      DTFloat32,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.GetDeprecationMessage(),
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Float32) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.Float32) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrNil(attr.Default, func(v defaults.Float32) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}
		case schema.Float64Attribute:
			field = ResourceField{
				Parents:       parents,
				Name:          name,
				DataType:      DTFloat64,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.DeprecationMessage,
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Float64) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.Float64) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrNil(attr.Default, func(v defaults.Float64) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}
		case schema.Int32Attribute:
			field = ResourceField{
				Parents:       parents,
				Name:          name,
				DataType:      DTInt32,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.GetDeprecationMessage(),
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Int32) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.Int32) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrNil(attr.Default, func(v defaults.Int32) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}
		case schema.Int64Attribute:
			field = ResourceField{
				Parents:       parents,
				Name:          name,
				DataType:      DTInt64,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.GetDeprecationMessage(),
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Int64) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.Int64) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrNil(attr.Default, func(v defaults.Int64) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}
		case schema.NumberAttribute:
			field = ResourceField{
				Parents:       parents,
				Name:          name,
				DataType:      DTNumber,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.GetDeprecationMessage(),
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Number) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.Number) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrNil(attr.Default, func(v defaults.Number) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}
		case schema.StringAttribute:
			field = ResourceField{
				Parents:       parents,
				Name:          name,
				DataType:      DTString,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.GetDeprecationMessage(),
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.String) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.String) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrNil(attr.Default, func(v defaults.String) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}
		case schema.ListAttribute:
			field = ResourceField{
				Parents:       parents,
				Name:          name,
				DataType:      DTList,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.GetDeprecationMessage(),
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.List) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.List) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrNil(attr.Default, func(v defaults.List) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}
		case schema.MapAttribute:
			field = ResourceField{
				Parents:       parents,
				Name:          name,
				DataType:      DTMap,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.GetDeprecationMessage(),
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Map) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.Map) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrNil(attr.Default, func(v defaults.Map) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}
		case schema.SetAttribute:
			field = ResourceField{
				Parents:       parents,
				Name:          name,
				DataType:      DTSet,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.GetDeprecationMessage(),
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Set) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.Set) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrNil(attr.Default, func(v defaults.Set) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}
		case schema.DynamicAttribute:
			field = ResourceField{
				Parents:       parents,
				Name:          name,
				DataType:      DTDynamic,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.GetDeprecationMessage(),
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Dynamic) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.Dynamic) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrNil(attr.Default, func(v defaults.Dynamic) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}

		case schema.ObjectAttribute:
			field = ResourceField{
				Parents:       parents,
				Name:          name,
				DataType:      DTObjectAttr,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.GetDeprecationMessage(),
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Object) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrNil(attr.Default, func(v defaults.Object) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}
			// NOTE: We don't look into the AttributeTypes for an ObjectAttribute as it doesn't contain useful information.
		case schema.SingleNestedAttribute:
			field = ResourceField{
				Parents:       parents,
				Name:          name,
				DataType:      DTSingleNestedAttr,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.GetDeprecationMessage(),
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Object) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrNil(attr.Default, func(v defaults.Object) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}
			objectNested, objectDiags = newResourceNestedAttrObjectFields(ctx, slices.Concat(parents, []string{name}), attr.GetNestedObject().(schema.NestedAttributeObject))
		case schema.SetNestedAttribute:
			field = ResourceField{
				Parents:       parents,
				Name:          name,
				DataType:      DTSetNestedAttr,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.GetDeprecationMessage(),
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Set) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.Set) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrNil(attr.Default, func(v defaults.Set) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}
			objectNested, objectDiags = newResourceNestedAttrObjectFields(ctx, slices.Concat(parents, []string{name}), attr.GetNestedObject().(schema.NestedAttributeObject))
		case schema.MapNestedAttribute:
			field = ResourceField{
				Parents:       parents,
				Name:          name,
				DataType:      DTMapNestedAttr,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.GetDeprecationMessage(),
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Map) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.Map) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrNil(attr.Default, func(v defaults.Map) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}
			objectNested, objectDiags = newResourceNestedAttrObjectFields(ctx, slices.Concat(parents, []string{name}), attr.GetNestedObject().(schema.NestedAttributeObject))
		case schema.ListNestedAttribute:
			field = ResourceField{
				Parents:       parents,
				Name:          name,
				DataType:      DTListNestedAttr,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.GetDeprecationMessage(),
				planModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.List) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(attr.Validators, func(v validator.List) string { return DescriptionCtxOf(ctx, v) }),
				defaultDesc:   MapOrNil(attr.Default, func(v defaults.List) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
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

func newResourceNestedAttrObjectFields(ctx context.Context, parents []string, obj schema.NestedAttributeObject) (nested ResourceNestedFields, diags diag.Diagnostics) {
	nested = ResourceNestedFields{}

	attrFields, attrNested, attrDiags := newResourceAttrFields(ctx, parents, obj.Attributes)
	diags.Append(attrDiags...)
	if diags.HasError() {
		return
	}

	nested[strings.Join(parents, ".")] = ResourceNestedField{
		PlanModifiers: MapSlice(obj.PlanModifiers, func(v planmodifier.Object) string { return DescriptionCtxOf(ctx, v) }),
		Validators:    MapSlice(obj.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
		Fields:        attrFields,
	}
	maps.Copy(nested, attrNested)
	return
}

func newResourceBlockFields(ctx context.Context, parents []string, blks map[string]schema.Block) (fields ResourceFields, nested ResourceNestedFields, diags diag.Diagnostics) {
	fields = ResourceFields{}
	nested = ResourceNestedFields{}

	for name, blk := range blks {
		var field ResourceField

		switch blk := blk.(type) {
		case schema.SingleNestedBlock:
			field = ResourceField{
				Parents:       parents,
				Name:          name,
				DataType:      DTSingleNestedBlock,
				Optional:      true, // Always regard a block as optional.
				Description:   DescriptionOf(blk),
				Deprecation:   blk.GetDeprecationMessage(),
				planModifiers: MapSlice(blk.PlanModifiers, func(v planmodifier.Object) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(blk.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.ListNestedBlock:
			field = ResourceField{
				Parents:       parents,
				Name:          name,
				DataType:      DTListNestedBlock,
				Optional:      true, // Always regard a block as optional.
				Description:   DescriptionOf(blk),
				Deprecation:   blk.GetDeprecationMessage(),
				planModifiers: MapSlice(blk.PlanModifiers, func(v planmodifier.List) string { return DescriptionCtxOf(ctx, v) }),
				validators:    MapSlice(blk.Validators, func(v validator.List) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.SetNestedBlock:
			field = ResourceField{
				Parents:       parents,
				Name:          name,
				DataType:      DTSetNestedBlock,
				Optional:      true, // Always regard a block as optional.
				Description:   DescriptionOf(blk),
				Deprecation:   blk.GetDeprecationMessage(),
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

func newResourceNestedBlkObjectFields(ctx context.Context, parents []string, obj schema.NestedBlockObject) (nested ResourceNestedFields, diags diag.Diagnostics) {
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

	fields := ResourceFields{}
	maps.Copy(fields, attrFields)
	maps.Copy(fields, blkFields)

	nested = ResourceNestedFields{}
	maps.Copy(nested, attrNested)
	maps.Copy(nested, blkNested)

	nested[strings.Join(parents, ".")] = ResourceNestedField{
		PlanModifiers: MapSlice(obj.PlanModifiers, func(v planmodifier.Object) string { return DescriptionCtxOf(ctx, v) }),
		Validators:    MapSlice(obj.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
		Fields:        fields,
	}
	return
}
