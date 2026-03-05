package metadata

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

type Fields map[string]Field

func (fields Fields) RequiredFields() []Field {
	var out []Field
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

func (fields Fields) OptionalFields() []Field {
	var out []Field
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

func (fields Fields) ComputedFields() []Field {
	var out []Field
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

type NestedFields map[string]NestedField

type NestedField struct {
	PlanModifiers []string
	Validators    []string
	Fields        Fields
}

type Field struct {
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

func (field Field) NestedKey() string {
	return strings.Join(slices.Concat(field.Parents, []string{field.Name}), ".")
}

func (field Field) NestedLink() string {
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

func (field Field) Traits() string {
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

func (field Field) Default() string {
	if field.defaultDesc == nil {
		return ""
	}
	return Sentencefy(*field.defaultDesc)
}

func (field Field) PlanModifiers() []string {
	var out []string
	for _, e := range field.planModifiers {
		out = append(out, Sentencefy(e))
	}
	return out
}

func (field Field) Validators() []string {
	var out []string
	for _, e := range field.validators {
		out = append(out, Sentencefy(e))
	}
	return out
}
