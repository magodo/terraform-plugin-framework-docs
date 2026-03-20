package fwdtypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ basetypes.StringTypable = StringType{}
var _ attr.TypeWithMarkdownDescription = StringType{}

type StringType struct {
	description string
	basetypes.StringType
}

func NewStringType(description string) StringType {
	return StringType{
		description: description,
		StringType:  basetypes.StringType{},
	}
}

// MarkdownDescription implements [attr.TypeWithMarkdownDescription].
func (s StringType) MarkdownDescription(context.Context) string {
	return s.description
}
