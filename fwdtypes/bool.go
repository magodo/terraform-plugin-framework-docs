package fwdtypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ basetypes.BoolTypable = BoolType{}
var _ attr.TypeWithMarkdownDescription = BoolType{}

type BoolType struct {
	description string
	basetypes.BoolType
}

func NewBoolType(description string) BoolType {
	return BoolType{
		description: description,
		BoolType:    basetypes.BoolType{},
	}
}

// MarkdownDescription implements [attr.TypeWithMarkdownDescription].
func (s BoolType) MarkdownDescription(context.Context) string {
	return s.description
}
