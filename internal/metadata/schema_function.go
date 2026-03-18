package metadata

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type FunctionSchema struct {
	Description string
	Summary     string
	Deprecation string

	Parameters FunctionFields
	Objects    FunctionObjects

	Return        FunctionField
	ReturnObjects FunctionObjects
}

func NewFunctionSchema(ctx context.Context, sch function.Definition) (schema FunctionSchema, diags diag.Diagnostics) {
	nested := FunctionObjects{}

	paramFields, paramNested, odiags := newFunctionParameters(ctx, nil, sch.Parameters)
	diags.Append(odiags...)
	if diags.HasError() {
		return
	}
	maps.Copy(nested, paramNested)

	if sch.VariadicParameter != nil {
		varFields, paramNested, odiags := newFunctionParameters(ctx, nil, []function.Parameter{sch.VariadicParameter})
		diags.Append(odiags...)
		if diags.HasError() {
			return
		}
		field := varFields[0]
		field.isVariadic = true

		paramFields = append(paramFields, field)
		maps.Copy(nested, paramNested)
	}

	retField, retNested, odiags := newFunctionReturn(ctx, sch.Return)
	diags.Append(odiags...)
	if diags.HasError() {
		return
	}

	schema = FunctionSchema{
		Summary: sch.Summary,
		Description: func() string {
			v := sch.MarkdownDescription
			if v != "" {
				return v
			}
			return sch.Description
		}(),
		Deprecation:   sch.DeprecationMessage,
		Parameters:    paramFields,
		Objects:       nested,
		Return:        retField,
		ReturnObjects: retNested,
	}
	return
}

func newFunctionReturn(ctx context.Context, ret function.Return) (field FunctionField, objects FunctionObjects, diags diag.Diagnostics) {
	switch attr := ret.(type) {
	case function.BoolReturn:
		field = FunctionField{
			dataType: DTBool,
		}
	case function.Float32Return:
		field = FunctionField{
			dataType: DTFloat32,
		}
	case function.Float64Return:
		field = FunctionField{
			dataType: DTFloat64,
		}
	case function.Int32Return:
		field = FunctionField{
			dataType: DTInt32,
		}
	case function.Int64Return:
		field = FunctionField{
			dataType: DTInt64,
		}
	case function.NumberReturn:
		field = FunctionField{
			dataType: DTNumber,
		}
	case function.StringReturn:
		field = FunctionField{
			dataType: DTString,
		}
	case function.ListReturn:
		field = FunctionField{
			dataType: DTList,
		}
	case function.MapReturn:
		field = FunctionField{
			dataType: DTMap,
		}
	case function.SetReturn:
		field = FunctionField{
			dataType: DTSet,
		}
	case function.DynamicReturn:
		field = FunctionField{
			dataType: DTDynamic,
		}
	case function.ObjectReturn:
		field = FunctionField{
			dataType: DTObjectAttr,
		}
		var odiags diag.Diagnostics
		objects, odiags = newFunctionObjects(ctx, nil, attr.AttributeTypes)
		diags = append(diags, odiags...)
		if diags.HasError() {
			return
		}
	default:
		diags.AddError("unknown schema type", fmt.Sprintf("%T", attr))
		return
	}

	return field, objects, diags
}

func newFunctionParameters(ctx context.Context, parents []string, params []function.Parameter) (fields FunctionFields, nested FunctionObjects, diags diag.Diagnostics) {
	fields = FunctionFields{}
	nested = FunctionObjects{}

	for _, attr := range params {
		var field FunctionField

		switch attr := attr.(type) {
		case function.BoolParameter:
			field = FunctionField{
				parents:      parents,
				name:         attr.GetName(),
				dataType:     DTBool,
				description:  DescriptionOf(attr),
				allowNull:    attr.GetAllowNullValue(),
				allowUnknown: attr.GetAllowUnknownValues(),
				validators: MapSliceSome(attr.Validators,
					func(v function.BoolParameterValidator) *string {
						return MaybeDescriptionCtxOf(ctx, v)
					}),
			}
		case function.Float32Parameter:
			field = FunctionField{
				parents:      parents,
				name:         attr.GetName(),
				dataType:     DTFloat32,
				description:  DescriptionOf(attr),
				allowNull:    attr.GetAllowNullValue(),
				allowUnknown: attr.GetAllowUnknownValues(),
				validators: MapSliceSome(attr.Validators,
					func(v function.Float32ParameterValidator) *string {
						return MaybeDescriptionCtxOf(ctx, v)
					}),
			}
		case function.Float64Parameter:
			field = FunctionField{
				parents:      parents,
				name:         attr.GetName(),
				dataType:     DTFloat64,
				description:  DescriptionOf(attr),
				allowNull:    attr.GetAllowNullValue(),
				allowUnknown: attr.GetAllowUnknownValues(),
				validators: MapSliceSome(attr.Validators,
					func(v function.Float64ParameterValidator) *string {
						return MaybeDescriptionCtxOf(ctx, v)
					}),
			}
		case function.Int32Parameter:
			field = FunctionField{
				parents:      parents,
				name:         attr.GetName(),
				dataType:     DTInt32,
				description:  DescriptionOf(attr),
				allowNull:    attr.GetAllowNullValue(),
				allowUnknown: attr.GetAllowUnknownValues(),
				validators: MapSliceSome(attr.Validators,
					func(v function.Int32ParameterValidator) *string {
						return MaybeDescriptionCtxOf(ctx, v)
					}),
			}
		case function.Int64Parameter:
			field = FunctionField{
				parents:      parents,
				name:         attr.GetName(),
				dataType:     DTInt64,
				description:  DescriptionOf(attr),
				allowNull:    attr.GetAllowNullValue(),
				allowUnknown: attr.GetAllowUnknownValues(),
				validators: MapSliceSome(attr.Validators,
					func(v function.Int64ParameterValidator) *string {
						return MaybeDescriptionCtxOf(ctx, v)
					}),
			}
		case function.NumberParameter:
			field = FunctionField{
				parents:      parents,
				name:         attr.GetName(),
				dataType:     DTNumber,
				description:  DescriptionOf(attr),
				allowNull:    attr.GetAllowNullValue(),
				allowUnknown: attr.GetAllowUnknownValues(),
				validators: MapSliceSome(attr.Validators,
					func(v function.NumberParameterValidator) *string {
						return MaybeDescriptionCtxOf(ctx, v)
					}),
			}
		case function.StringParameter:
			field = FunctionField{
				parents:      parents,
				name:         attr.GetName(),
				dataType:     DTString,
				description:  DescriptionOf(attr),
				allowNull:    attr.GetAllowNullValue(),
				allowUnknown: attr.GetAllowUnknownValues(),
				validators: MapSliceSome(attr.Validators,
					func(v function.StringParameterValidator) *string {
						return MaybeDescriptionCtxOf(ctx, v)
					}),
			}
		case function.ListParameter:
			field = FunctionField{
				parents:      parents,
				name:         attr.GetName(),
				dataType:     DTList,
				description:  DescriptionOf(attr),
				allowNull:    attr.GetAllowNullValue(),
				allowUnknown: attr.GetAllowUnknownValues(),
				validators: MapSliceSome(attr.Validators,
					func(v function.ListParameterValidator) *string {
						return MaybeDescriptionCtxOf(ctx, v)
					}),
			}
		case function.MapParameter:
			field = FunctionField{
				parents:      parents,
				name:         attr.GetName(),
				dataType:     DTMap,
				description:  DescriptionOf(attr),
				allowNull:    attr.GetAllowNullValue(),
				allowUnknown: attr.GetAllowUnknownValues(),
				validators: MapSliceSome(attr.Validators,
					func(v function.MapParameterValidator) *string {
						return MaybeDescriptionCtxOf(ctx, v)
					}),
			}
		case function.SetParameter:
			field = FunctionField{
				parents:      parents,
				name:         attr.GetName(),
				dataType:     DTSet,
				description:  DescriptionOf(attr),
				allowNull:    attr.GetAllowNullValue(),
				allowUnknown: attr.GetAllowUnknownValues(),
				validators: MapSliceSome(attr.Validators,
					func(v function.SetParameterValidator) *string {
						return MaybeDescriptionCtxOf(ctx, v)
					}),
			}
		case function.DynamicParameter:
			field = FunctionField{
				parents:      parents,
				name:         attr.GetName(),
				dataType:     DTDynamic,
				description:  DescriptionOf(attr),
				allowNull:    attr.GetAllowNullValue(),
				allowUnknown: attr.GetAllowUnknownValues(),
				validators: MapSliceSome(attr.Validators,
					func(v function.DynamicParameterValidator) *string {
						return MaybeDescriptionCtxOf(ctx, v)
					}),
			}
		case function.ObjectParameter:
			field = FunctionField{
				parents:      parents,
				name:         attr.GetName(),
				dataType:     DTObjectAttr,
				description:  DescriptionOf(attr),
				allowNull:    attr.GetAllowNullValue(),
				allowUnknown: attr.GetAllowUnknownValues(),
				validators: MapSliceSome(attr.Validators,
					func(v function.ObjectParameterValidator) *string {
						return MaybeDescriptionCtxOf(ctx, v)
					}),
			}
			nestedObjects, odiags := newFunctionObjects(ctx, slices.Concat(parents, []string{attr.GetName()}), attr.AttributeTypes)
			diags = append(diags, odiags...)
			if diags.HasError() {
				return
			}
			maps.Copy(nested, nestedObjects)

		default:
			diags.AddError("unknown schema type", fmt.Sprintf("%T", attr))
			return
		}

		fields = append(fields, field)
	}

	return
}

func newFunctionObjects(ctx context.Context, parents []string, attrs map[string]attr.Type) (objects FunctionObjects, diags diag.Diagnostics) {
	objects = FunctionObjects{}
	fields := map[string]FunctionField{}
	for name, attr := range attrs {
		var field FunctionField
		switch attr := attr.(type) {
		case basetypes.BoolType:
			field = FunctionField{
				parents:     parents,
				name:        name,
				dataType:    DTBool,
				description: PointerTo(MaybeDescriptionCtxOf(ctx, attr)),
			}
		case basetypes.Float32Type:
			field = FunctionField{
				parents:     parents,
				name:        name,
				dataType:    DTFloat32,
				description: PointerTo(MaybeDescriptionCtxOf(ctx, attr)),
			}
		case basetypes.Float64Type:
			field = FunctionField{
				parents:     parents,
				name:        name,
				dataType:    DTFloat64,
				description: PointerTo(MaybeDescriptionCtxOf(ctx, attr)),
			}
		case basetypes.Int32Type:
			field = FunctionField{
				parents:     parents,
				name:        name,
				dataType:    DTInt32,
				description: PointerTo(MaybeDescriptionCtxOf(ctx, attr)),
			}
		case basetypes.Int64Type:
			field = FunctionField{
				parents:     parents,
				name:        name,
				dataType:    DTInt64,
				description: PointerTo(MaybeDescriptionCtxOf(ctx, attr)),
			}
		case basetypes.NumberType:
			field = FunctionField{
				parents:     parents,
				name:        name,
				dataType:    DTNumber,
				description: PointerTo(MaybeDescriptionCtxOf(ctx, attr)),
			}
		case basetypes.StringType:
			field = FunctionField{
				parents:     parents,
				name:        name,
				dataType:    DTString,
				description: PointerTo(MaybeDescriptionCtxOf(ctx, attr)),
			}
		case basetypes.ListType:
			field = FunctionField{
				parents:     parents,
				name:        name,
				dataType:    DTList,
				description: PointerTo(MaybeDescriptionCtxOf(ctx, attr)),
			}
		case basetypes.SetType:
			field = FunctionField{
				parents:     parents,
				name:        name,
				dataType:    DTSet,
				description: PointerTo(MaybeDescriptionCtxOf(ctx, attr)),
			}
		case basetypes.MapType:
			field = FunctionField{
				parents:     parents,
				name:        name,
				dataType:    DTMap,
				description: PointerTo(MaybeDescriptionCtxOf(ctx, attr)),
			}
		case basetypes.DynamicType:
			field = FunctionField{
				parents:     parents,
				name:        name,
				dataType:    DTDynamic,
				description: PointerTo(MaybeDescriptionCtxOf(ctx, attr)),
			}
		case basetypes.ObjectType:
			field = FunctionField{
				parents:     parents,
				name:        name,
				dataType:    DTObjectAttr,
				description: PointerTo(MaybeDescriptionCtxOf(ctx, attr)),
			}
			nestedObjects, odiags := newFunctionObjects(ctx, slices.Concat(parents, []string{name}), attr.AttributeTypes())
			diags = append(diags, odiags...)
			if diags.HasError() {
				return nil, diags
			}
			maps.Copy(objects, nestedObjects)
		}

		fields[name] = field
	}
	objects[strings.Join(parents, ".")] = FunctionObject{
		functionKey: strings.Join(parents, "."),
		fields:      fields,
	}

	return objects, diags
}
