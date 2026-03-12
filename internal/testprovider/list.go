package testprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type ExampleList struct{}

var _ list.ListResource = ExampleList{}

// Metadata implements [list.ListResource].
func (e ExampleList) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource"
}

// ListResourceConfigSchema implements [list.ListResource].
func (e ExampleList) ListResourceConfigSchema(ctx context.Context, req list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
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
		MarkdownDescription: "List example resources.",
		Attributes: map[string]schema.Attribute{
			"bool": schema.BoolAttribute{
				MarkdownDescription: "A boolean attribute.",
				DeprecationMessage:  "Deprecated in favor of `boolean`.",
				Required:            true,
				Validators: []validator.Bool{
					boolvalidator.AlsoRequires(path.MatchRoot("string"), path.MatchRoot("int64")),
					boolvalidator.ConflictsWith(path.MatchRoot("list")),
				},
			},
			"string": schema.StringAttribute{
				MarkdownDescription: "A string attribute.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("foo", "bar", "baz"),
				},
			},
			"int64": schema.Int64Attribute{
				MarkdownDescription: "A int64 attribute.",
				Optional:            true,
			},
			"list": schema.ListAttribute{
				MarkdownDescription: "A list attribute.",
				Optional:            true,
			},
			"map": schema.MapAttribute{
				MarkdownDescription: "A map attribute.",
				Optional:            true,
			},
			"dynamic": schema.DynamicAttribute{
				MarkdownDescription: "A dynamic attribute.",
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
		},
	}
}

// List implements [list.ListResource].
func (e ExampleList) List(context.Context, list.ListRequest, *list.ListResultsStream) {
	panic("unimplemented")
}
