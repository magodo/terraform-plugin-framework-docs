package metadata

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type DataSourceSchema struct {
	Description string
	Deprecation string

	Fields DataSourceFields

	// Including nested attribute object or block object.
	Nested DataSourceNestedFields
}

func NewDataSourceSchema(ctx context.Context, sch schema.Schema) (schema DataSourceSchema, diags diag.Diagnostics) {
	fields := DataSourceFields{}
	nested := DataSourceNestedFields{}

	attrFields, attrNested, odiags := newDataSourceAttrFields(ctx, nil, sch.Attributes)
	diags.Append(odiags...)
	if diags.HasError() {
		return
	}
	maps.Copy(fields, attrFields)
	maps.Copy(nested, attrNested)

	blockFields, blockNested, odiags := newDataSourceBlockFields(ctx, nil, sch.Blocks)
	diags.Append(odiags...)
	if diags.HasError() {
		return
	}
	maps.Copy(fields, blockFields)
	maps.Copy(nested, blockNested)

	schema = DataSourceSchema{
		Description: DescriptionOf(sch),
		Deprecation: sch.GetDeprecationMessage(),
		Fields:      fields,
		Nested:      nested,
	}
	return
}

func newDataSourceAttrFields(ctx context.Context, parents []string, attrs map[string]schema.Attribute) (fields DataSourceFields, nested DataSourceNestedFields, diags diag.Diagnostics) {
	fields = DataSourceFields{}
	nested = DataSourceNestedFields{}

	for name, attr := range attrs {
		var (
			field DataSourceField

			objectNested DataSourceNestedFields
			objectDiags  diag.Diagnostics
		)

		switch attr := attr.(type) {
		case schema.BoolAttribute:
			field = DataSourceField{
				Parents:     parents,
				Name:        name,
				DataType:    DTBool,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Computed:    attr.IsComputed(),
				Sensitive:   attr.IsSensitive(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Bool) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.Float32Attribute:
			field = DataSourceField{
				Parents:     parents,
				Name:        name,
				DataType:    DTFloat32,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Computed:    attr.IsComputed(),
				Sensitive:   attr.IsSensitive(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Float32) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.Float64Attribute:
			field = DataSourceField{
				Parents:     parents,
				Name:        name,
				DataType:    DTFloat64,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Computed:    attr.IsComputed(),
				Sensitive:   attr.IsSensitive(),
				Description: DescriptionOf(attr),
				Deprecation: attr.DeprecationMessage,
				validators:  MapSlice(attr.Validators, func(v validator.Float64) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.Int32Attribute:
			field = DataSourceField{
				Parents:     parents,
				Name:        name,
				DataType:    DTInt32,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Computed:    attr.IsComputed(),
				Sensitive:   attr.IsSensitive(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Int32) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.Int64Attribute:
			field = DataSourceField{
				Parents:     parents,
				Name:        name,
				DataType:    DTInt64,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Computed:    attr.IsComputed(),
				Sensitive:   attr.IsSensitive(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Int64) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.NumberAttribute:
			field = DataSourceField{
				Parents:     parents,
				Name:        name,
				DataType:    DTNumber,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Computed:    attr.IsComputed(),
				Sensitive:   attr.IsSensitive(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Number) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.StringAttribute:
			field = DataSourceField{
				Parents:     parents,
				Name:        name,
				DataType:    DTString,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Computed:    attr.IsComputed(),
				Sensitive:   attr.IsSensitive(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.String) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.ListAttribute:
			field = DataSourceField{
				Parents:     parents,
				Name:        name,
				DataType:    DTList,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Computed:    attr.IsComputed(),
				Sensitive:   attr.IsSensitive(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.List) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.MapAttribute:
			field = DataSourceField{
				Parents:     parents,
				Name:        name,
				DataType:    DTMap,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Computed:    attr.IsComputed(),
				Sensitive:   attr.IsSensitive(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Map) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.SetAttribute:
			field = DataSourceField{
				Parents:     parents,
				Name:        name,
				DataType:    DTSet,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Computed:    attr.IsComputed(),
				Sensitive:   attr.IsSensitive(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Set) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.DynamicAttribute:
			field = DataSourceField{
				Parents:     parents,
				Name:        name,
				DataType:    DTDynamic,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Computed:    attr.IsComputed(),
				Sensitive:   attr.IsSensitive(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Dynamic) string { return DescriptionCtxOf(ctx, v) }),
			}

		case schema.ObjectAttribute:
			field = DataSourceField{
				Parents:     parents,
				Name:        name,
				DataType:    DTObjectAttr,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Computed:    attr.IsComputed(),
				Sensitive:   attr.IsSensitive(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
			}
			// NOTE: We don't look into the AttributeTypes for an ObjectAttribute as it doesn't contain useful information.
		case schema.SingleNestedAttribute:
			field = DataSourceField{
				Parents:     parents,
				Name:        name,
				DataType:    DTSingleNestedAttr,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Computed:    attr.IsComputed(),
				Sensitive:   attr.IsSensitive(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
			}
			objectNested, objectDiags = newDataSourceNestedAttrObjectFields(ctx, slices.Concat(parents, []string{name}), attr.GetNestedObject().(schema.NestedAttributeObject))
		case schema.SetNestedAttribute:
			field = DataSourceField{
				Parents:     parents,
				Name:        name,
				DataType:    DTSetNestedAttr,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Computed:    attr.IsComputed(),
				Sensitive:   attr.IsSensitive(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Set) string { return DescriptionCtxOf(ctx, v) }),
			}
			objectNested, objectDiags = newDataSourceNestedAttrObjectFields(ctx, slices.Concat(parents, []string{name}), attr.GetNestedObject().(schema.NestedAttributeObject))
		case schema.MapNestedAttribute:
			field = DataSourceField{
				Parents:     parents,
				Name:        name,
				DataType:    DTMapNestedAttr,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Computed:    attr.IsComputed(),
				Sensitive:   attr.IsSensitive(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.Map) string { return DescriptionCtxOf(ctx, v) }),
			}
			objectNested, objectDiags = newDataSourceNestedAttrObjectFields(ctx, slices.Concat(parents, []string{name}), attr.GetNestedObject().(schema.NestedAttributeObject))
		case schema.ListNestedAttribute:
			field = DataSourceField{
				Parents:     parents,
				Name:        name,
				DataType:    DTListNestedAttr,
				Required:    attr.IsRequired(),
				Optional:    attr.IsOptional(),
				Computed:    attr.IsComputed(),
				Sensitive:   attr.IsSensitive(),
				Description: DescriptionOf(attr),
				Deprecation: attr.GetDeprecationMessage(),
				validators:  MapSlice(attr.Validators, func(v validator.List) string { return DescriptionCtxOf(ctx, v) }),
			}
			objectNested, objectDiags = newDataSourceNestedAttrObjectFields(ctx, slices.Concat(parents, []string{name}), attr.GetNestedObject().(schema.NestedAttributeObject))
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

func newDataSourceNestedAttrObjectFields(ctx context.Context, parents []string, obj schema.NestedAttributeObject) (nested DataSourceNestedFields, diags diag.Diagnostics) {
	nested = DataSourceNestedFields{}

	attrFields, attrNested, attrDiags := newDataSourceAttrFields(ctx, parents, obj.Attributes)
	diags.Append(attrDiags...)
	if diags.HasError() {
		return
	}

	nested[strings.Join(parents, ".")] = DataSourceNestedField{
		Validators: MapSlice(obj.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
		Fields:     attrFields,
	}
	maps.Copy(nested, attrNested)
	return
}

func newDataSourceBlockFields(ctx context.Context, parents []string, blks map[string]schema.Block) (fields DataSourceFields, nested DataSourceNestedFields, diags diag.Diagnostics) {
	fields = DataSourceFields{}
	nested = DataSourceNestedFields{}

	for name, blk := range blks {
		var field DataSourceField

		switch blk := blk.(type) {
		case schema.SingleNestedBlock:
			field = DataSourceField{
				Parents:     parents,
				Name:        name,
				DataType:    DTSingleNestedBlock,
				Optional:    true, // Always regard a block as optional.
				Description: DescriptionOf(blk),
				Deprecation: blk.GetDeprecationMessage(),
				validators:  MapSlice(blk.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.ListNestedBlock:
			field = DataSourceField{
				Parents:     parents,
				Name:        name,
				DataType:    DTListNestedBlock,
				Optional:    true, // Always regard a block as optional.
				Description: DescriptionOf(blk),
				Deprecation: blk.GetDeprecationMessage(),
				validators:  MapSlice(blk.Validators, func(v validator.List) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.SetNestedBlock:
			field = DataSourceField{
				Parents:     parents,
				Name:        name,
				DataType:    DTSetNestedBlock,
				Optional:    true, // Always regard a block as optional.
				Description: DescriptionOf(blk),
				Deprecation: blk.GetDeprecationMessage(),
				validators:  MapSlice(blk.Validators, func(v validator.Set) string { return DescriptionCtxOf(ctx, v) }),
			}
		}

		objectNested, odiags := newDataSourceNestedBlkObjectFields(ctx, slices.Concat(parents, []string{name}), blk.GetNestedObject().(schema.NestedBlockObject))
		diags = append(diags, odiags...)
		if diags.HasError() {
			return
		}

		fields[name] = field
		maps.Copy(nested, objectNested)
	}

	return
}

func newDataSourceNestedBlkObjectFields(ctx context.Context, parents []string, obj schema.NestedBlockObject) (nested DataSourceNestedFields, diags diag.Diagnostics) {
	attrFields, attrNested, attrDiags := newDataSourceAttrFields(ctx, parents, obj.Attributes)
	diags.Append(attrDiags...)
	if diags.HasError() {
		return
	}

	blkFields, blkNested, attrDiags := newDataSourceBlockFields(ctx, parents, obj.Blocks)
	diags.Append(attrDiags...)
	if diags.HasError() {
		return
	}

	fields := DataSourceFields{}
	maps.Copy(fields, attrFields)
	maps.Copy(fields, blkFields)

	nested = DataSourceNestedFields{}
	maps.Copy(nested, attrNested)
	maps.Copy(nested, blkNested)

	nested[strings.Join(parents, ".")] = DataSourceNestedField{
		Validators: MapSlice(obj.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
		Fields:     fields,
	}
	return
}
