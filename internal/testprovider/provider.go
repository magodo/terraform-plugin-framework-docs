package testprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &ExampleCloudProvider{}
var _ provider.ProviderWithEphemeralResources = &ExampleCloudProvider{}
var _ provider.ProviderWithActions = &ExampleCloudProvider{}
var _ provider.ProviderWithListResources = &ExampleCloudProvider{}
var _ provider.ProviderWithFunctions = &ExampleCloudProvider{}

type ExampleCloudProvider struct{}

func (p *ExampleCloudProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "examplecloud"
}

func (p *ExampleCloudProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	nestedAttrs := map[string]schema.Attribute{
		"bool": schema.BoolAttribute{
			MarkdownDescription: "A nested bool attribute.",
			DeprecationMessage:  "Deprecated in favor of `boolean`.",
			Required:            true,
		},
		"string": schema.StringAttribute{
			MarkdownDescription: "A nested string attribute.",
			Optional:            true,
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
		MarkdownDescription: "The example provider.",
		DeprecationMessage:  "This provider is deprecated.",
		Attributes: map[string]schema.Attribute{
			"bool": schema.BoolAttribute{
				MarkdownDescription: "A boolean attribute.",
				DeprecationMessage:  "Deprecated in favor of `boolean`.",
				Required:            true,
				Sensitive:           true,
				Validators: []validator.Bool{
					boolvalidator.AlsoRequires(path.MatchRoot("string"), path.MatchRoot("int64")),
					boolvalidator.ConflictsWith(path.MatchRoot("list")),
				},
			},
			"string": schema.StringAttribute{
				MarkdownDescription: "A string attribute.",
				Optional:            true,
				Sensitive:           true,
				Validators: []validator.String{
					stringvalidator.OneOf("foo", "bar", "baz"),
				},
			},
			"int64": schema.Int64Attribute{
				MarkdownDescription: "A int64 attribute.",
				Optional:            true,
			},
			"list": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "A list attribute.",
				Optional:            true,
			},
			"map": schema.MapAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "A map attribute.",
				Optional:            true,
			},
			"set": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "A set attribute.",
				Optional:            true,
			},
			"dynamic": schema.DynamicAttribute{
				MarkdownDescription: "A dynamic attribute.",
				Optional:            true,
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
		},
		Blocks: map[string]schema.Block{
			"single_block": schema.SingleNestedBlock{
				MarkdownDescription: "A single block.",
				Attributes:          nestedAttrs,
				Blocks:              nestedBlks,
				Validators: []validator.Object{
					objectvalidator.ConflictsWith(path.MatchRoot("list_block")),
				},
			},
			"list_block": schema.ListNestedBlock{
				MarkdownDescription: "A list block.",
				NestedObject: schema.NestedBlockObject{
					Attributes: nestedAttrs,
					Blocks:     nestedBlks,
					Validators: []validator.Object{
						objectvalidator.IsRequired(),
					},
				},
			},
			"set_block": schema.SetNestedBlock{
				MarkdownDescription: "A set block.",
				NestedObject: schema.NestedBlockObject{
					Attributes: nestedAttrs,
					Blocks:     nestedBlks,
				},
			},
		},
	}
}

func (p *ExampleCloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

func (p *ExampleCloudProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return ExampleResource{} },
	}
}

func (p *ExampleCloudProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource { return ExampleDataSource{} },
	}
}

func (p *ExampleCloudProvider) EphemeralResources(context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		func() ephemeral.EphemeralResource { return ExampleEphemeralResource{} },
	}
}

func (p *ExampleCloudProvider) Actions(context.Context) []func() action.Action {
	return []func() action.Action{
		func() action.Action { return ExampleAction{} },
	}
}

func (p *ExampleCloudProvider) ListResources(context.Context) []func() list.ListResource {
	return []func() list.ListResource{
		func() list.ListResource { return ExampleList{} },
	}
}

func (p *ExampleCloudProvider) Functions(context.Context) []func() function.Function {
	return []func() function.Function{
		func() function.Function { return ExampleFunctionSimple{} },
		func() function.Function { return ExampleFunctionRetObj{} },
	}
}
