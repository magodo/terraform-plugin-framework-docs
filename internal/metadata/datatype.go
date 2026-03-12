package metadata

import "fmt"

// Terraform data types.
// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types
type DataType int

const (
	// Primary
	DTBool DataType = iota
	DTFloat32
	DTFloat64
	DTInt32
	DTInt64
	DTNumber
	DTString
	// Collection
	DTList
	DTMap
	DTSet
	// Object (attribute)
	DTSingleNestedAttr
	DTListNestedAttr
	DTMapNestedAttr
	DTSetNestedAttr
	DTObjectAttr
	// Object (block)
	DTSingleNestedBlock
	DTListNestedBlock
	DTSetNestedBlock
	// Tuple
	DTTuple
	// Dynamic
	DTDynamic
)

func (dt DataType) String() string {
	switch dt {
	case DTBool:
		return "Boolean"
	case DTFloat32:
		return "Float32"
	case DTFloat64:
		return "Float64"
	case DTInt32:
		return "Int32"
	case DTInt64:
		return "Int64"
	case DTNumber:
		return "Number"
	case DTString:
		return "String"
	case DTList:
		return "List"
	case DTMap:
		return "Map"
	case DTSet:
		return "Set"
	case DTSingleNestedAttr:
		return "Single Object"
	case DTListNestedAttr:
		return "List of Objects"
	case DTMapNestedAttr:
		return "Map of Objects"
	case DTSetNestedAttr:
		return "Set of Objects"
	case DTObjectAttr:
		return "Object"
	case DTSingleNestedBlock:
		return "Single Block"
	case DTListNestedBlock, DTSetNestedBlock:
		return "Blocks"
	case DTTuple:
		return "Tuple"
	case DTDynamic:
		return "Dynamic"
	default:
		panic(fmt.Sprintf("unknown data type %d", dt))
	}
}
