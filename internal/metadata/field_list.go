package metadata

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

type ListFields map[string]ListField

func (fields ListFields) RequiredFields() []ListField {
	var out []ListField
	for _, info := range fields {
		if !info.Required {
			continue
		}
		out = append(out, info)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})
	return out
}

func (fields ListFields) OptionalFields() []ListField {
	var out []ListField
	for _, info := range fields {
		if !info.Optional {
			continue
		}
		out = append(out, info)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})
	return out
}

type ListNestedFields map[string]ListNestedField

type ListNestedField struct {
	PlanModifiers []string
	Validators    []string
	Fields        ListFields
}

type ListField struct {
	Parents  []string
	Name     string
	DataType DataType

	Required bool
	Optional bool

	Description string
	Deprecation string

	validators []string
}

func (field ListField) NestedKey() string {
	return strings.Join(slices.Concat(field.Parents, []string{field.Name}), ".")
}

func (field ListField) NestedLink() string {
	switch field.DataType {
	case DTSingleNestedAttr,
		DTListNestedAttr,
		DTMapNestedAttr,
		DTSetNestedAttr,
		DTObjectAttr,
		DTSingleNestedBlock,
		DTListNestedBlock,
		DTSetNestedBlock:
		return fmt.Sprintf("See the nested schema [here](#nested--%s).", field.NestedKey())
	default:
		return ""
	}
}

func (field ListField) Traits() string {
	var traits []string
	traits = append(traits, field.DataType.String())
	return strings.Join(traits, ", ")
}

func (field ListField) Validators() []string {
	var out []string
	for _, e := range field.validators {
		out = append(out, Sentencefy(e))
	}
	return out
}
