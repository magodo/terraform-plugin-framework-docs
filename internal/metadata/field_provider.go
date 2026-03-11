package metadata

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

type ProviderFields map[string]ProviderField

func (fields ProviderFields) RequiredFields() []ProviderField {
	var out []ProviderField
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

func (fields ProviderFields) OptionalFields() []ProviderField {
	var out []ProviderField
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

type ProviderNestedFields map[string]ProviderNestedField

type ProviderNestedField struct {
	Validators []string
	Fields     ProviderFields
}

type ProviderField struct {
	Parents  []string
	Name     string
	DataType DataType

	Required bool
	Optional bool

	Sensitive bool

	Description string
	Deprecation string

	validators []string
}

func (field ProviderField) NestedKey() string {
	return strings.Join(slices.Concat(field.Parents, []string{field.Name}), ".")
}

func (field ProviderField) NestedLink() string {
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

func (field ProviderField) Traits() string {
	var traits []string
	traits = append(traits, field.DataType.String())
	if field.Sensitive {
		traits = append(traits, "Sensitive")
	}
	return strings.Join(traits, ", ")
}

func (field ProviderField) Validators() []string {
	var out []string
	for _, e := range field.validators {
		out = append(out, Sentencefy(e))
	}
	return out
}
