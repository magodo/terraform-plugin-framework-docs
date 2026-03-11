package metadata

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

type EphemeralFields map[string]EphemeralField

func (fields EphemeralFields) RequiredFields() []EphemeralField {
	var out []EphemeralField
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

func (fields EphemeralFields) OptionalFields() []EphemeralField {
	var out []EphemeralField
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

func (fields EphemeralFields) ComputedFields() []EphemeralField {
	var out []EphemeralField
	for _, info := range fields {
		if !(info.Computed && !info.Optional) {
			continue
		}
		out = append(out, info)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})
	return out
}

type EphemeralNestedFields map[string]EphemeralNestedField

type EphemeralNestedField struct {
	Validators []string
	Fields     EphemeralFields
}

type EphemeralField struct {
	Parents  []string
	Name     string
	DataType DataType

	Required bool
	Optional bool
	Computed bool

	Sensitive bool

	Description string
	Deprecation string

	validators []string
}

func (field EphemeralField) NestedKey() string {
	return strings.Join(slices.Concat(field.Parents, []string{field.Name}), ".")
}

func (field EphemeralField) NestedLink() string {
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

func (field EphemeralField) Traits() string {
	var traits []string
	traits = append(traits, field.DataType.String())
	if field.Sensitive {
		traits = append(traits, "Sensitive")
	}
	return strings.Join(traits, ", ")
}

func (field EphemeralField) Validators() []string {
	var out []string
	for _, e := range field.validators {
		out = append(out, Sentencefy(e))
	}
	return out
}
