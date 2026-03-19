package metadata

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/magodo/terraform-plugin-framework-docs/internal/testhelper"
	"github.com/stretchr/testify/require"
)

func TestDataType(t *testing.T) {
	cases := []struct {
		name   string
		dt     DataType
		expect string
	}{
		{
			name: "bool",
			dt: DataType{
				inner: basetypes.BoolType{},
			},
			expect: "Boolean",
		},
		{
			name: "float32",
			dt: DataType{
				inner: basetypes.Float32Type{},
			},
			expect: "Float32",
		},
		{
			name: "float64",
			dt: DataType{
				inner: basetypes.Float64Type{},
			},
			expect: "Float64",
		},
		{
			name: "int32",
			dt: DataType{
				inner: basetypes.Int32Type{},
			},
			expect: "Int32",
		},
		{
			name: "int64",
			dt: DataType{
				inner: basetypes.Int64Type{},
			},
			expect: "Int64",
		},
		{
			name: "number",
			dt: DataType{
				inner: basetypes.NumberType{},
			},
			expect: "Number",
		},
		{
			name: "string",
			dt: DataType{
				inner: basetypes.StringType{},
			},
			expect: "String",
		},
		{
			name: "dynamic",
			dt: DataType{
				inner: basetypes.DynamicType{},
			},
			expect: "Dynamic",
		},
		{
			name: "single object",
			dt: DataType{
				inner: basetypes.ObjectType{},
			},
			expect: "Object",
		},
		{
			name: "list of objects",
			dt: DataType{
				inner: basetypes.ListType{
					ElemType: basetypes.ObjectType{},
				},
			},
			expect: "List of Objects",
		},
		{
			name: "set of objects",
			dt: DataType{
				inner: basetypes.SetType{
					ElemType: basetypes.ObjectType{},
				},
			},
			expect: "Set of Objects",
		},
		{
			name: "map of objects",
			dt: DataType{
				inner: basetypes.MapType{
					ElemType: basetypes.ObjectType{},
				},
			},
			expect: "Map of Objects",
		},
		{
			name: "list of blocks",
			dt: DataType{
				isblk: true,
				inner: basetypes.ListType{
					ElemType: basetypes.ObjectType{},
				},
			},
			expect: "List of Blocks",
		},
		{
			name: "set of blocks",
			dt: DataType{
				isblk: true,
				inner: basetypes.SetType{
					ElemType: basetypes.ObjectType{},
				},
			},
			expect: "Set of Blocks",
		},
		{
			name: "map of blocks",
			dt: DataType{
				isblk: true,
				inner: basetypes.MapType{
					ElemType: basetypes.ObjectType{},
				},
			},
			expect: "Map of Blocks",
		},
		{
			name: "tuple",
			dt: DataType{
				inner: basetypes.TupleType{},
			},
			expect: "Tuple",
		},
		{
			name: "custom string",
			dt: DataType{
				inner: testhelper.CustomStringType{},
			},
			expect: "String",
		},
		{
			name: "list of maps of objects",
			dt: DataType{
				inner: basetypes.ListType{
					ElemType: basetypes.MapType{
						ElemType: basetypes.ObjectType{},
					},
				},
			},
			expect: "List of Maps of Objects",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expect, tt.dt.String())
		})
	}
}
