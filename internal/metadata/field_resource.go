package metadata

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

type ResourceFields map[string]ResourceField

func (fields ResourceFields) RequiredFields() []ResourceField {
	var out []ResourceField
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

func (fields ResourceFields) OptionalFields() []ResourceField {
	var out []ResourceField
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

func (fields ResourceFields) ComputedFields() []ResourceField {
	var out []ResourceField
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

type ResourceNestedFields map[string]ResourceNestedField

type ResourceNestedField struct {
	PlanModifiers []string
	Validators    []string
	Fields        ResourceFields
}

type ResourceField struct {
	Parents  []string
	Name     string
	DataType DataType

	Required bool
	Optional bool
	Computed bool

	Sensitive bool

	Description string
	Deprecation string

	WriteOnly bool

	defaultDesc   *string
	planModifiers []string
	validators    []string
}

func (field ResourceField) NestedKey() string {
	return strings.Join(slices.Concat(field.Parents, []string{field.Name}), ".")
}

func (field ResourceField) NestedLink() string {
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

func (field ResourceField) Traits() string {
	var traits []string
	traits = append(traits, field.DataType.String())
	if field.Sensitive {
		traits = append(traits, "Sensitive")
	}
	if field.WriteOnly {
		traits = append(traits, "[Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)")
	}
	return strings.Join(traits, ", ")
}

func (field ResourceField) Default() string {
	if field.defaultDesc == nil {
		return ""
	}
	return Sentencefy(*field.defaultDesc)
}

func (field ResourceField) PlanModifiers() []string {
	var out []string
	for _, e := range field.planModifiers {
		out = append(out, Sentencefy(e))
	}
	return out
}

func (field ResourceField) Validators() []string {
	var out []string
	for _, e := range field.validators {
		out = append(out, Sentencefy(e))
	}
	return out
}
