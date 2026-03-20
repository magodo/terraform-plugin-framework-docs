package fwdtypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ basetypes.SetTypable = SetType{}
var _ attr.TypeWithMarkdownDescription = SetType{}

type SetType struct {
	description string
	basetypes.SetType
}

func NewSetType(description string, elemType attr.Type) SetType {
	return SetType{
		description: description,
		SetType:     basetypes.SetType{ElemType: elemType},
	}
}

// MarkdownDescription implements [attr.TypeWithMarkdownDescription].
func (s SetType) MarkdownDescription(context.Context) string {
	return s.description
}
