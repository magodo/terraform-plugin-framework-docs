package testprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ExampleFunctionRetObj struct{}

var _ function.Function = ExampleFunctionRetObj{}

// Metadata implements [function.Function].
func (e ExampleFunctionRetObj) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "example_function_retobj"
}

// Definition implements [function.Function].
func (e ExampleFunctionRetObj) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name: "object",
				AttributeTypes: map[string]attr.Type{
					"foo": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"bar": types.BoolType,
						},
					},
				},
				MarkdownDescription: "An object parameter.",
			},
		},
		Return: function.ObjectReturn{
			AttributeTypes: map[string]attr.Type{
				"retfoo": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"retbar": types.BoolType,
					},
				},
			},
		},
		Summary:             "The summary.",
		MarkdownDescription: "The description.",
	}
}

// Run implements [function.Function].
func (e ExampleFunctionRetObj) Run(context.Context, function.RunRequest, *function.RunResponse) {
	panic("unimplemented")
}
