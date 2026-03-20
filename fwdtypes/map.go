package fwdtypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ basetypes.ObjectTypable = ObjectType{}
var _ attr.TypeWithMarkdownDescription = ObjectType{}

type ObjectType struct {
	description string
	basetypes.ObjectType
}

func NewObjectType(description string, attrTypes map[string]attr.Type) ObjectType {
	return ObjectType{
		description: description,
		ObjectType:  basetypes.ObjectType{AttrTypes: attrTypes},
	}
}

// MarkdownDescription implements [attr.TypeWithMarkdownDescription].
func (s ObjectType) MarkdownDescription(context.Context) string {
	return s.description
}
