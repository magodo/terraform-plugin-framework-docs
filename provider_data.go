package tfproviderdocs

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type ProviderData struct {
	Resources ResourceInfos
}

func NewProviderData(ctx context.Context, p provider.Provider) (pd ProviderData, diags diag.Diagnostics) {
	pd = ProviderData{
		Resources: ResourceInfos{},
	}

	var providerMetadataResp provider.MetadataResponse
	p.Metadata(ctx, provider.MetadataRequest{}, &providerMetadataResp)

	for _, builder := range p.Resources(ctx) {
		res := builder()

		// Get the resource type
		var metadataResp resource.MetadataResponse
		res.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: providerMetadataResp.TypeName}, &metadataResp)
		resourceType := metadataResp.TypeName

		var schemaResp resource.SchemaResponse
		res.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
		diags.Append(schemaResp.Diagnostics...)
		if diags.HasError() {
			return
		}

		info, odiags := pd.NewResourceInfo(ctx, schemaResp.Schema)
		diags.Append(odiags...)
		if diags.HasError() {
			return
		}

		pd.Resources[resourceType] = info
	}

	return
}

func (pd ProviderData) NewResourceInfo(ctx context.Context, sch schema.Schema) (resourceInfo ResourceInfo, diags diag.Diagnostics) {
	infos := SchemaInfos{}
	nested := NestedSchemaInfos{}

	attrInfos, attrNested, odiags := pd.newResourceAttrInfos(ctx, nil, sch.Attributes)
	diags.Append(odiags...)
	if diags.HasError() {
		return
	}
	maps.Copy(infos, attrInfos)
	maps.Copy(nested, attrNested)

	blockInfos, blockNested, odiags := pd.newResourceBlockInfos(ctx, nil, sch.Blocks)
	maps.Copy(infos, blockInfos)
	maps.Copy(nested, blockNested)

	resourceInfo = ResourceInfo{
		Description: DescriptionOf(sch),
		Deprecation: sch.GetDeprecationMessage(),
		Infos:       infos,
		Nested:      nested,
	}
	return
}

func (pd ProviderData) newResourceAttrInfos(ctx context.Context, parents []string, attrs map[string]schema.Attribute) (infos SchemaInfos, nested NestedSchemaInfos, diags diag.Diagnostics) {
	infos = SchemaInfos{}
	nested = NestedSchemaInfos{}

	for name, attr := range attrs {
		var (
			info SchemaInfo

			objectNested NestedSchemaInfos
			objectDiags  diag.Diagnostics
		)

		switch attr := attr.(type) {
		case schema.BoolAttribute:
			info = SchemaInfo{
				Parents:       parents,
				Name:          name,
				DataType:      DTBool,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.GetDeprecationMessage(),
				PlanModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Bool) string { return DescriptionCtxOf(ctx, v) }),
				Validators:    MapSlice(attr.Validators, func(v validator.Bool) string { return DescriptionCtxOf(ctx, v) }),
				DefaultDesc:   MapOrNil(attr.Default, func(v defaults.Bool) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}
		case schema.Float32Attribute:
			info = SchemaInfo{
				Parents:       parents,
				Name:          name,
				DataType:      DTFloat32,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.GetDeprecationMessage(),
				PlanModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Float32) string { return DescriptionCtxOf(ctx, v) }),
				Validators:    MapSlice(attr.Validators, func(v validator.Float32) string { return DescriptionCtxOf(ctx, v) }),
				DefaultDesc:   MapOrNil(attr.Default, func(v defaults.Float32) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}
		case schema.Float64Attribute:
			info = SchemaInfo{
				Parents:       parents,
				Name:          name,
				DataType:      DTFloat64,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.DeprecationMessage,
				PlanModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Float64) string { return DescriptionCtxOf(ctx, v) }),
				Validators:    MapSlice(attr.Validators, func(v validator.Float64) string { return DescriptionCtxOf(ctx, v) }),
				DefaultDesc:   MapOrNil(attr.Default, func(v defaults.Float64) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}
		case schema.Int32Attribute:
			info = SchemaInfo{
				Parents:       parents,
				Name:          name,
				DataType:      DTInt32,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.GetDeprecationMessage(),
				PlanModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Int32) string { return DescriptionCtxOf(ctx, v) }),
				Validators:    MapSlice(attr.Validators, func(v validator.Int32) string { return DescriptionCtxOf(ctx, v) }),
				DefaultDesc:   MapOrNil(attr.Default, func(v defaults.Int32) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}
		case schema.Int64Attribute:
			info = SchemaInfo{
				Parents:       parents,
				Name:          name,
				DataType:      DTInt64,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.GetDeprecationMessage(),
				PlanModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Int64) string { return DescriptionCtxOf(ctx, v) }),
				Validators:    MapSlice(attr.Validators, func(v validator.Int64) string { return DescriptionCtxOf(ctx, v) }),
				DefaultDesc:   MapOrNil(attr.Default, func(v defaults.Int64) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}
		case schema.NumberAttribute:
			info = SchemaInfo{
				Parents:       parents,
				Name:          name,
				DataType:      DTNumber,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.GetDeprecationMessage(),
				PlanModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Number) string { return DescriptionCtxOf(ctx, v) }),
				Validators:    MapSlice(attr.Validators, func(v validator.Number) string { return DescriptionCtxOf(ctx, v) }),
				DefaultDesc:   MapOrNil(attr.Default, func(v defaults.Number) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}
		case schema.StringAttribute:
			info = SchemaInfo{
				Parents:       parents,
				Name:          name,
				DataType:      DTString,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.GetDeprecationMessage(),
				PlanModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.String) string { return DescriptionCtxOf(ctx, v) }),
				Validators:    MapSlice(attr.Validators, func(v validator.String) string { return DescriptionCtxOf(ctx, v) }),
				DefaultDesc:   MapOrNil(attr.Default, func(v defaults.String) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}
		case schema.ListAttribute:
			info = SchemaInfo{
				Parents:       parents,
				Name:          name,
				DataType:      DTList,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.GetDeprecationMessage(),
				PlanModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.List) string { return DescriptionCtxOf(ctx, v) }),
				Validators:    MapSlice(attr.Validators, func(v validator.List) string { return DescriptionCtxOf(ctx, v) }),
				DefaultDesc:   MapOrNil(attr.Default, func(v defaults.List) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}
		case schema.MapAttribute:
			info = SchemaInfo{
				Parents:       parents,
				Name:          name,
				DataType:      DTMap,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.GetDeprecationMessage(),
				PlanModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Map) string { return DescriptionCtxOf(ctx, v) }),
				Validators:    MapSlice(attr.Validators, func(v validator.Map) string { return DescriptionCtxOf(ctx, v) }),
				DefaultDesc:   MapOrNil(attr.Default, func(v defaults.Map) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}
		case schema.SetAttribute:
			info = SchemaInfo{
				Parents:       parents,
				Name:          name,
				DataType:      DTSet,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.GetDeprecationMessage(),
				PlanModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Set) string { return DescriptionCtxOf(ctx, v) }),
				Validators:    MapSlice(attr.Validators, func(v validator.Set) string { return DescriptionCtxOf(ctx, v) }),
				DefaultDesc:   MapOrNil(attr.Default, func(v defaults.Set) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}
		case schema.DynamicAttribute:
			info = SchemaInfo{
				Parents:       parents,
				Name:          name,
				DataType:      DTDynamic,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.GetDeprecationMessage(),
				PlanModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Dynamic) string { return DescriptionCtxOf(ctx, v) }),
				Validators:    MapSlice(attr.Validators, func(v validator.Dynamic) string { return DescriptionCtxOf(ctx, v) }),
				DefaultDesc:   MapOrNil(attr.Default, func(v defaults.Dynamic) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}

		case schema.ObjectAttribute:
			info = SchemaInfo{
				Parents:       parents,
				Name:          name,
				DataType:      DTObjectAttr,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.GetDeprecationMessage(),
				PlanModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Object) string { return DescriptionCtxOf(ctx, v) }),
				Validators:    MapSlice(attr.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
				DefaultDesc:   MapOrNil(attr.Default, func(v defaults.Object) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}
			// NOTE: We don't look into the AttributeTypes for an ObjectAttribute as it doesn't contain useful information.
		case schema.SingleNestedAttribute:
			info = SchemaInfo{
				Parents:       parents,
				Name:          name,
				DataType:      DTSingleNestedAttr,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.GetDeprecationMessage(),
				PlanModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Object) string { return DescriptionCtxOf(ctx, v) }),
				Validators:    MapSlice(attr.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
				DefaultDesc:   MapOrNil(attr.Default, func(v defaults.Object) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}
			objectNested, objectDiags = pd.newResourceNestedAttrObjectInfos(ctx, slices.Concat(parents, []string{name}), attr.GetNestedObject().(schema.NestedAttributeObject))
		case schema.SetNestedAttribute:
			info = SchemaInfo{
				Parents:       parents,
				Name:          name,
				DataType:      DTSetNestedAttr,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.GetDeprecationMessage(),
				PlanModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Set) string { return DescriptionCtxOf(ctx, v) }),
				Validators:    MapSlice(attr.Validators, func(v validator.Set) string { return DescriptionCtxOf(ctx, v) }),
				DefaultDesc:   MapOrNil(attr.Default, func(v defaults.Set) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}
			objectNested, objectDiags = pd.newResourceNestedAttrObjectInfos(ctx, slices.Concat(parents, []string{name}), attr.GetNestedObject().(schema.NestedAttributeObject))
		case schema.MapNestedAttribute:
			info = SchemaInfo{
				Parents:       parents,
				Name:          name,
				DataType:      DTMapNestedAttr,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.GetDeprecationMessage(),
				PlanModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.Map) string { return DescriptionCtxOf(ctx, v) }),
				Validators:    MapSlice(attr.Validators, func(v validator.Map) string { return DescriptionCtxOf(ctx, v) }),
				DefaultDesc:   MapOrNil(attr.Default, func(v defaults.Map) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}
			objectNested, objectDiags = pd.newResourceNestedAttrObjectInfos(ctx, slices.Concat(parents, []string{name}), attr.GetNestedObject().(schema.NestedAttributeObject))
		case schema.ListNestedAttribute:
			info = SchemaInfo{
				Parents:       parents,
				Name:          name,
				DataType:      DTListNestedAttr,
				Required:      attr.IsRequired(),
				Optional:      attr.IsOptional(),
				Computed:      attr.IsComputed(),
				Sensitive:     attr.IsSensitive(),
				Description:   DescriptionOf(attr),
				Deprecation:   attr.GetDeprecationMessage(),
				PlanModifiers: MapSlice(attr.PlanModifiers, func(v planmodifier.List) string { return DescriptionCtxOf(ctx, v) }),
				Validators:    MapSlice(attr.Validators, func(v validator.List) string { return DescriptionCtxOf(ctx, v) }),
				DefaultDesc:   MapOrNil(attr.Default, func(v defaults.List) string { return DescriptionCtxOf(ctx, v) }),
				WriteOnly:     attr.IsWriteOnly(),
			}
			objectNested, objectDiags = pd.newResourceNestedAttrObjectInfos(ctx, slices.Concat(parents, []string{name}), attr.GetNestedObject().(schema.NestedAttributeObject))
		default:
			diags.AddError("unknown schema type", fmt.Sprintf("%T", attr))
			return
		}

		infos[name] = info

		diags = append(diags, objectDiags...)
		if diags.HasError() {
			return
		}
		maps.Copy(nested, objectNested)
	}

	return
}

func (pd ProviderData) newResourceNestedAttrObjectInfos(ctx context.Context, parents []string, obj schema.NestedAttributeObject) (nested NestedSchemaInfos, diags diag.Diagnostics) {
	nested = NestedSchemaInfos{}

	attrInfos, attrNested, attrDiags := pd.newResourceAttrInfos(ctx, parents, obj.Attributes)
	diags.Append(attrDiags...)
	if diags.HasError() {
		return
	}

	nested[strings.Join(parents, ".")] = NestedSchemaInfo{
		PlanModifiers: MapSlice(obj.PlanModifiers, func(v planmodifier.Object) string { return DescriptionCtxOf(ctx, v) }),
		Validators:    MapSlice(obj.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
		Infos:         attrInfos,
	}
	maps.Copy(nested, attrNested)
	return
}

func (pd ProviderData) newResourceBlockInfos(ctx context.Context, parents []string, blks map[string]schema.Block) (infos SchemaInfos, nested NestedSchemaInfos, diags diag.Diagnostics) {
	infos = SchemaInfos{}
	nested = NestedSchemaInfos{}

	for name, blk := range blks {
		var info SchemaInfo

		switch blk := blk.(type) {
		case schema.SingleNestedBlock:
			info = SchemaInfo{
				Parents:       parents,
				Name:          name,
				DataType:      DTSingleNestedBlock,
				Optional:      true, // Always regard a block as optional.
				Description:   DescriptionOf(blk),
				Deprecation:   blk.GetDeprecationMessage(),
				PlanModifiers: MapSlice(blk.PlanModifiers, func(v planmodifier.Object) string { return DescriptionCtxOf(ctx, v) }),
				Validators:    MapSlice(blk.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.ListNestedBlock:
			info = SchemaInfo{
				Parents:       parents,
				Name:          name,
				DataType:      DTListNestedBlock,
				Optional:      true, // Always regard a block as optional.
				Description:   DescriptionOf(blk),
				Deprecation:   blk.GetDeprecationMessage(),
				PlanModifiers: MapSlice(blk.PlanModifiers, func(v planmodifier.List) string { return DescriptionCtxOf(ctx, v) }),
				Validators:    MapSlice(blk.Validators, func(v validator.List) string { return DescriptionCtxOf(ctx, v) }),
			}
		case schema.SetNestedBlock:
			info = SchemaInfo{
				Parents:       parents,
				Name:          name,
				DataType:      DTSetNestedBlock,
				Optional:      true, // Always regard a block as optional.
				Description:   DescriptionOf(blk),
				Deprecation:   blk.GetDeprecationMessage(),
				PlanModifiers: MapSlice(blk.PlanModifiers, func(v planmodifier.Set) string { return DescriptionCtxOf(ctx, v) }),
				Validators:    MapSlice(blk.Validators, func(v validator.Set) string { return DescriptionCtxOf(ctx, v) }),
			}
		}

		objectNested, odiags := pd.newResourceNestedBlkObjectInfos(ctx, slices.Concat(parents, []string{name}), blk.GetNestedObject().(schema.NestedBlockObject))
		diags = append(diags, odiags...)
		if diags.HasError() {
			return
		}

		infos[name] = info
		maps.Copy(nested, objectNested)
	}

	return
}

func (pd ProviderData) newResourceNestedBlkObjectInfos(ctx context.Context, parents []string, obj schema.NestedBlockObject) (nested NestedSchemaInfos, diags diag.Diagnostics) {
	attrInfos, attrNested, attrDiags := pd.newResourceAttrInfos(ctx, parents, obj.Attributes)
	diags.Append(attrDiags...)
	if diags.HasError() {
		return
	}

	blkInfos, blkNested, attrDiags := pd.newResourceBlockInfos(ctx, parents, obj.Blocks)
	diags.Append(attrDiags...)
	if diags.HasError() {
		return
	}

	infos := SchemaInfos{}
	maps.Copy(infos, attrInfos)
	maps.Copy(infos, blkInfos)

	nested = NestedSchemaInfos{}
	maps.Copy(nested, attrNested)
	maps.Copy(nested, blkNested)

	nested[strings.Join(parents, ".")] = NestedSchemaInfo{
		PlanModifiers: MapSlice(obj.PlanModifiers, func(v planmodifier.Object) string { return DescriptionCtxOf(ctx, v) }),
		Validators:    MapSlice(obj.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
		Infos:         infos,
	}
	return
}
