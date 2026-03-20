package fwdtypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ basetypes.MapTypable = MapType{}
var _ attr.TypeWithMarkdownDescription = MapType{}

type MapType struct {
	description string
	basetypes.MapType
}

func NewMapType(description string, elemType attr.Type) MapType {
	return MapType{
		description: description,
		MapType:     basetypes.MapType{ElemType: elemType},
	}
}

// MarkdownDescription implements [attr.TypeWithMarkdownDescription].
func (s MapType) MarkdownDescription(context.Context) string {
	return s.description
}
