package testprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/magodo/terraform-plugin-framework-docs/fwdtypes"
	"github.com/magodo/terraform-plugin-framework-docs/internal/testhelper"
)

type ExampleFunctionSimple struct{}

var _ function.Function = ExampleFunctionSimple{}

// Metadata implements [function.Function].
func (e ExampleFunctionSimple) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "example_function_simple"
}

// Definition implements [function.Function].
func (e ExampleFunctionSimple) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Parameters: []function.Parameter{
			function.BoolParameter{
				Name:                "bool",
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				MarkdownDescription: "A bool parameter.",
				Validators: []function.BoolParameterValidator{
					boolvalidator.Equals(true),
				},
			},
			function.StringParameter{
				Name:                "string",
				MarkdownDescription: "A string parameter",
			},
			function.NumberParameter{
				Name:                "number",
				MarkdownDescription: "A number parameter",
			},
			function.Int32Parameter{
				Name:                "int32",
				MarkdownDescription: "A int32 parameter",
			},
			function.Int64Parameter{
				Name:                "int64",
				MarkdownDescription: "A int64 parameter",
			},
			function.Float32Parameter{
				Name:                "float32",
				MarkdownDescription: "A float32 parameter",
			},
			function.Float64Parameter{
				Name:                "float64",
				MarkdownDescription: "A float64 parameter",
			},
			function.DynamicParameter{
				Name:                "dynamic",
				MarkdownDescription: "A dynamic parameter",
			},
			function.ListParameter{
				Name:                "list",
				ElementType:         types.StringType,
				MarkdownDescription: "A list of string parameter.",
			},
			function.SetParameter{
				Name:                "set",
				ElementType:         types.StringType,
				MarkdownDescription: "A set of string parameter.",
			},
			function.MapParameter{
				Name:                "map",
				ElementType:         types.StringType,
				MarkdownDescription: "A map of string parameter.",
			},
			function.ObjectParameter{
				Name: "object",
				AttributeTypes: map[string]attr.Type{
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
					"object": fwdtypes.NewObjectType(
						"Description.",
						map[string]attr.Type{
							"foo": fwdtypes.NewBoolType("A foo field."),
						},
					),
				},
				MarkdownDescription: "An object parameter.",
			},
			function.StringParameter{
				Name:                "custom_string",
				MarkdownDescription: "A custom string attribute.",
				CustomType:          testhelper.CustomStringType{},
			},
		},
		VariadicParameter: function.StringParameter{
			Name:                "strings",
			MarkdownDescription: "The variadic string parameter.",
		},
		Return:              function.BoolReturn{},
		Summary:             "The summary.",
		MarkdownDescription: "The description.",
		DeprecationMessage:  "This is deprecated.",
	}
}

// Run implements [function.Function].
func (e ExampleFunctionSimple) Run(context.Context, function.RunRequest, *function.RunResponse) {
	panic("unimplemented")
}
