package fwdtypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ basetypes.ListTypable = ListType{}
var _ attr.TypeWithMarkdownDescription = ListType{}

type ListType struct {
	description string
	basetypes.ListType
}

func NewListType(description string, elemType attr.Type) ListType {
	return ListType{
		description: description,
		ListType:    basetypes.ListType{ElemType: elemType},
	}
}

// MarkdownDescription implements [attr.TypeWithMarkdownDescription].
func (s ListType) MarkdownDescription(context.Context) string {
	return s.description
}
