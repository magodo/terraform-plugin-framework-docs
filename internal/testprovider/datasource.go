package testprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ExampleDataSource struct{}

var _ datasource.DataSource = ExampleDataSource{}

func (e ExampleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource"
}

// Schema implements [datasource.DataSource].
func (e ExampleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
		MarkdownDescription: "Queries an example resource.",
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
				Computed:            true,
				Sensitive:           true,
				Validators: []validator.String{
					stringvalidator.OneOf("foo", "bar", "baz"),
				},
			},
			"int64": schema.Int64Attribute{
				MarkdownDescription: "A int64 attribute.",
				Optional:            true,
				Computed:            true,
			},
			"list": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "A list attribute.",
				Optional:            true,
				Computed:            true,
			},
			"map": schema.MapAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "A map attribute.",
				Optional:            true,
				Computed:            true,
			},
			"set": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "A set attribute.",
				Optional:            true,
				Computed:            true,
			},
			"dynamic": schema.DynamicAttribute{
				MarkdownDescription: "A dynamic attribute.",
				Computed:            true,
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

// Read implements [datasource.DataSource].
func (e ExampleDataSource) Read(context.Context, datasource.ReadRequest, *datasource.ReadResponse) {
	panic("unimplemented")
}
