package testprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/magodo/terraform-plugin-framework-docs/fwdtypes"
	"github.com/magodo/terraform-plugin-framework-docs/internal/testhelper"
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
			CustomType: testhelper.CustomObjectType{},
			AttributeTypes: map[string]attr.Type{
				"retfoo": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"retbar":  types.BoolType,
						"bool":    fwdtypes.NewBoolType("Description"),
						"int32":   fwdtypes.NewInt32Type("Description"),
						"int64":   fwdtypes.NewInt64Type("Description"),
						"float32": fwdtypes.NewFloat32Type("Description"),
						"float64": fwdtypes.NewFloat64Type("Description"),
						"number":  fwdtypes.NewNumberType("Description"),
						"string":  fwdtypes.NewStringType("Description"),
						"dynamic": fwdtypes.NewDynamicType("Description"),
						"list":    fwdtypes.NewListType("Description", basetypes.BoolType{}),
						"set":     fwdtypes.NewSetType("Description", basetypes.BoolType{}),
						"map":     fwdtypes.NewMapType("Description", basetypes.BoolType{}),
						"object":  fwdtypes.NewObjectType("Description", nil),
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
