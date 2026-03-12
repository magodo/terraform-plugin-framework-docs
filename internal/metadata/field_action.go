package metadata

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

type ActionFields map[string]ActionField

func (fields ActionFields) RequiredFields() []ActionField {
	var out []ActionField
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

func (fields ActionFields) OptionalFields() []ActionField {
	var out []ActionField
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

type ActionNestedFields map[string]ActionNestedField

type ActionNestedField struct {
	Validators []string
	Fields     ActionFields
}

type ActionField struct {
	Parents  []string
	Name     string
	DataType DataType

	Required bool
	Optional bool

	Description string
	Deprecation string

	WriteOnly bool

	validators []string
}

func (field ActionField) NestedKey() string {
	return strings.Join(slices.Concat(field.Parents, []string{field.Name}), ".")
}

func (field ActionField) NestedLink() string {
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

func (field ActionField) Traits() string {
	var traits []string
	traits = append(traits, field.DataType.String())
	if field.WriteOnly {
		traits = append(traits, "[Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)")
	}
	return strings.Join(traits, ", ")
}

func (field ActionField) Validators() []string {
	var out []string
	for _, e := range field.validators {
		out = append(out, Sentencefy(e))
	}
	return out
}
