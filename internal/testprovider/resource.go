package testprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/magodo/terraform-plugin-framework-docs/fwdtypes"
	"github.com/magodo/terraform-plugin-framework-docs/internal/testhelper"
)

type ExampleResource struct{}

var _ resource.ResourceWithIdentity = ExampleResource{}

// Metadata implements [resource.Resource].
func (e ExampleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource"
}

// Schema implements [resource.Resource].
func (e ExampleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	nestedAttrs := map[string]schema.Attribute{
		"bool": schema.BoolAttribute{
			MarkdownDescription: "A nested bool attribute.",
			DeprecationMessage:  "Deprecated in favor of `boolean`.",
			Required:            true,
		},
		"string": schema.StringAttribute{
			MarkdownDescription: "A nested string attribute.",
			Optional:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf("foo", "bar", "baz"),
			},
		},
		"nested_object": schema.SingleNestedAttribute{
			MarkdownDescription: "A nested single object attribute.",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"bool": schema.BoolAttribute{
					MarkdownDescription: "A nested nested bool attribute.",
					Required:            true,
				},
				"string": schema.StringAttribute{
					MarkdownDescription: "A nested nested string attribute.",
					Optional:            true,
				},
			},
		},
	}

	nestedBlks := map[string]schema.Block{
		"nested_block": schema.SingleNestedBlock{
			MarkdownDescription: "A nested block.",
			Attributes: map[string]schema.Attribute{
				"number": schema.NumberAttribute{
					MarkdownDescription: "A nested number attribute.",
					Optional:            true,
				},
			},
		},
	}
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an example resource.",
		Attributes: map[string]schema.Attribute{
			"bool": schema.BoolAttribute{
				MarkdownDescription: "A boolean attribute.",
				DeprecationMessage:  "Deprecated in favor of `boolean`.",
				Required:            true,
				WriteOnly:           true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
					boolplanmodifier.UseStateForUnknown(),
					boolplanmodifier.RequiresReplaceIf(func(ctx context.Context, br planmodifier.BoolRequest, rrifr *boolplanmodifier.RequiresReplaceIfFuncResponse) {
					}, "", "A conditional requires replace if."),
				},
				Validators: []validator.Bool{
					boolvalidator.AlsoRequires(path.MatchRoot("string"), path.MatchRoot("int64")),
					boolvalidator.ConflictsWith(path.MatchRoot("list")),
				},
			},
			"string": schema.StringAttribute{
				MarkdownDescription: "A string attribute.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				WriteOnly:           true,
				Default:             stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.OneOf("foo", "bar", "baz"),
				},
			},
			"int64": schema.Int64Attribute{
				MarkdownDescription: "A int64 attribute.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
			},
			"list": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "A list attribute.",
				Optional:            true,
				Computed:            true,
				Default:             listdefault.StaticValue(basetypes.NewListValueMust(basetypes.StringType{}, []attr.Value{basetypes.NewStringValue("foo")})),
			},
			"map": schema.MapAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "A map attribute.",
				Optional:            true,
				Computed:            true,
				Default:             mapdefault.StaticValue(basetypes.NewMapValueMust(basetypes.StringType{}, map[string]attr.Value{"key": basetypes.NewStringValue("val")})),
			},
			"set": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "A set attribute.",
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(basetypes.NewSetValueMust(basetypes.StringType{}, []attr.Value{basetypes.NewStringValue("foo")})),
			},
			"dynamic": schema.DynamicAttribute{
				MarkdownDescription: "A dynamic attribute.",
				Computed:            true,
			},
			"object": schema.ObjectAttribute{
				MarkdownDescription: "An object attribute.",
				Optional:            true,
				AttributeTypes: map[string]attr.Type{
					"foo": fwdtypes.NewObjectType(
						"A foo field.",
						map[string]attr.Type{
							"bool":    fwdtypes.NewBoolType("Description."),
							"int32":   fwdtypes.NewInt32Type("Description."),
							"int64":   fwdtypes.NewInt64Type("Description."),
							"float32": fwdtypes.NewFloat32Type("Description."),
							"float64": fwdtypes.NewFloat64Type("Description."),
							"number":  fwdtypes.NewNumberType("Description."),
							"string":  fwdtypes.NewStringType("Description."),
							"dynamic": fwdtypes.NewDynamicType("Description."),
							"list":    fwdtypes.NewListType("Description.", basetypes.BoolType{}),
							"set":     fwdtypes.NewSetType("Description.", basetypes.BoolType{}),
							"map":     fwdtypes.NewMapType("Description.", basetypes.BoolType{}),
						},
					),
				},
			},
			"single_object": schema.SingleNestedAttribute{
				MarkdownDescription: "A single object attribute.",
				Optional:            true,
				Attributes:          nestedAttrs,
			},
			"list_object": schema.ListNestedAttribute{
				MarkdownDescription: "A list object attribute.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					PlanModifiers: []planmodifier.Object{
						objectplanmodifier.RequiresReplace(),
					},
					Attributes: nestedAttrs,
				},
			},
			"map_object": schema.MapNestedAttribute{
				MarkdownDescription: "A map object attribute.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: nestedAttrs,
				},
			},
			"set_object": schema.SetNestedAttribute{
				MarkdownDescription: "A set object attribute.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: nestedAttrs,
				},
			},
			"custom_string": schema.StringAttribute{
				MarkdownDescription: "A custom string attribute.",
				CustomType:          testhelper.CustomStringType{},
				Optional:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"single_block": schema.SingleNestedBlock{
				MarkdownDescription: "A single block.",
				Attributes:          nestedAttrs,
				Blocks:              nestedBlks,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Object{
					objectvalidator.ConflictsWith(path.MatchRoot("list_block")),
				},
			},
			"list_block": schema.ListNestedBlock{
				MarkdownDescription: "A list block.",
				NestedObject: schema.NestedBlockObject{
					Attributes: nestedAttrs,
					Blocks:     nestedBlks,
					PlanModifiers: []planmodifier.Object{
						objectplanmodifier.RequiresReplace(),
					},
					Validators: []validator.Object{
						objectvalidator.IsRequired(),
					},
				},
				Validators: []validator.List{
					listvalidator.AlsoRequires(path.MatchRoot("single_block")),
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"set_block": schema.SetNestedBlock{
				MarkdownDescription: "A set block.",
				NestedObject: schema.NestedBlockObject{
					Attributes: nestedAttrs,
					Blocks:     nestedBlks,
					PlanModifiers: []planmodifier.Object{
						objectplanmodifier.RequiresReplace(),
					},
				},
			},
			"custom_block": schema.SingleNestedBlock{
				CustomType:          testhelper.CustomObjectType{},
				MarkdownDescription: "A custom block.",
				Attributes: map[string]schema.Attribute{
					"foo": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "A foo attribute.",
					},
				},
			},
		},
	}
}

// IdentitySchema implements [resource.ResourceWithIdentity].
func (e ExampleResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"parent_id": identityschema.StringAttribute{
				Description:       "The parent id.",
				RequiredForImport: true,
			},
			"id": identityschema.StringAttribute{
				Description:       "The id of this resource.",
				RequiredForImport: true,
			},
			"version": identityschema.StringAttribute{
				Description:       "The version of this resource.",
				OptionalForImport: true,
			},
			"custom_string": identityschema.StringAttribute{
				Description:       "A custom string attribute.",
				CustomType:        testhelper.CustomStringType{},
				OptionalForImport: true,
			},
		},
	}
}

// Create implements [resource.Resource].
func (e ExampleResource) Create(context.Context, resource.CreateRequest, *resource.CreateResponse) {
	panic("unimplemented")
}

// Delete implements [resource.Resource].
func (e ExampleResource) Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse) {
	panic("unimplemented")
}

// Read implements [resource.Resource].
func (e ExampleResource) Read(context.Context, resource.ReadRequest, *resource.ReadResponse) {
	panic("unimplemented")
}

// Update implements [resource.Resource].
func (e ExampleResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {
	panic("unimplemented")
}
