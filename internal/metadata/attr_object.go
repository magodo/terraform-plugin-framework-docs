package metadata

import (
	"context"
	"maps"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type Objects map[string]Object

type Object struct {
	key    string
	fields map[string]ObjectField
}

type ObjectField struct {
	parents     []string
	name        string
	dataType    DataType
	description string
}

func newObjects(ctx context.Context, parents []string, attrs map[string]attr.Type) (objects Objects, diags diag.Diagnostics) {
	objects = Objects{}
	fields := map[string]ObjectField{}
	for name, attr := range attrs {
		var field ObjectField
		switch attr := attr.(type) {
		case basetypes.BoolType:
			field = ObjectField{
				name:        name,
				dataType:    DTBool,
				description: PointerTo(MaybeDescriptionCtxOf(ctx, attr)),
			}
		case basetypes.Float32Type:
			field = ObjectField{
				parents:     parents,
				name:        name,
				dataType:    DTFloat32,
				description: PointerTo(MaybeDescriptionCtxOf(ctx, attr)),
			}
		case basetypes.Float64Type:
			field = ObjectField{
				parents:     parents,
				name:        name,
				dataType:    DTFloat64,
				description: PointerTo(MaybeDescriptionCtxOf(ctx, attr)),
			}
		case basetypes.Int32Type:
			field = ObjectField{
				parents:     parents,
				name:        name,
				dataType:    DTInt32,
				description: PointerTo(MaybeDescriptionCtxOf(ctx, attr)),
			}
		case basetypes.Int64Type:
			field = ObjectField{
				parents:     parents,
				name:        name,
				dataType:    DTInt64,
				description: PointerTo(MaybeDescriptionCtxOf(ctx, attr)),
			}
		case basetypes.NumberType:
			field = ObjectField{
				parents:     parents,
				name:        name,
				dataType:    DTNumber,
				description: PointerTo(MaybeDescriptionCtxOf(ctx, attr)),
			}
		case basetypes.StringType:
			field = ObjectField{
				parents:     parents,
				name:        name,
				dataType:    DTString,
				description: PointerTo(MaybeDescriptionCtxOf(ctx, attr)),
			}
		case basetypes.ListType:
			field = ObjectField{
				parents:     parents,
				name:        name,
				dataType:    DTList,
				description: PointerTo(MaybeDescriptionCtxOf(ctx, attr)),
			}
		case basetypes.SetType:
			field = ObjectField{
				parents:     parents,
				name:        name,
				dataType:    DTSet,
				description: PointerTo(MaybeDescriptionCtxOf(ctx, attr)),
			}
		case basetypes.MapType:
			field = ObjectField{
				parents:     parents,
				name:        name,
				dataType:    DTMap,
				description: PointerTo(MaybeDescriptionCtxOf(ctx, attr)),
			}
		case basetypes.DynamicType:
			field = ObjectField{
				parents:     parents,
				name:        name,
				dataType:    DTDynamic,
				description: PointerTo(MaybeDescriptionCtxOf(ctx, attr)),
			}
		case basetypes.ObjectType:
			field = ObjectField{
				parents:     parents,
				name:        name,
				dataType:    DTObjectAttr,
				description: PointerTo(MaybeDescriptionCtxOf(ctx, attr)),
			}
			nestedObjects, odiags := newObjects(ctx, slices.Concat(parents, []string{name}), attr.AttributeTypes())
			diags = append(diags, odiags...)
			if diags.HasError() {
				return nil, diags
			}
			maps.Copy(objects, nestedObjects)
		}

		fields[name] = field
	}
	key := strings.Join(parents, ".")
	objects[key] = Object{
		key:    key,
		fields: fields,
	}

	return objects, diags
}

func (objs Objects) ToFunctionObjects() FunctionObjects {
	out := FunctionObjects{}
	for key, obj := range objs {
		fields := map[string]FunctionField{}
		for key, field := range obj.fields {
			fields[key] = field.ToFunctionField()
		}
		out[key] = FunctionObject{
			functionKey: obj.key,
			fields:      fields,
		}
	}
	return out
}

func (obj ObjectField) ToFunctionField() FunctionField {
	return FunctionField{
		parents:     obj.parents,
		name:        obj.name,
		dataType:    obj.dataType,
		description: obj.description,
	}
}

func (objs Objects) ToNestedFields(rootField Field) NestedFields {
	out := NestedFields{}
	for key, obj := range objs {
		fields := Fields{}
		for key, field := range obj.fields {
			fields[key] = field.ToField(rootField)
		}
		out[key] = NestedField{
			fields: fields,
		}
	}
	return out
}

func (obj ObjectField) ToField(rootField Field) Field {
	return Field{
		parents:     obj.parents,
		name:        obj.name,
		dataType:    obj.dataType,
		required:    rootField.required,
		optional:    rootField.optional,
		computed:    rootField.computed,
		description: obj.description,
	}
}
