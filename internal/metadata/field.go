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
	for _, field := range fields {
		if !field.Required() {
			continue
		}
		out = append(out, field)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Name() < out[j].Name()
	})
	return out
}

func (fields Fields) OptionalFields() []Field {
	var out []Field
	for _, field := range fields {
		if !field.Optional() {
			continue
		}
		out = append(out, field)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Name() < out[j].Name()
	})
	return out
}

func (fields Fields) ComputedFields() []Field {
	var out []Field
	for _, field := range fields {
		if !(field.Computed() && !field.Optional()) {
			continue
		}
		out = append(out, field)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Name() < out[j].Name()
	})
	return out
}

type NestedFields map[string]NestedField

type NestedField struct {
	planModifiers []string
	validators    []string
	fields        Fields
}

func (r NestedField) Fields() Fields {
	fields := Fields{}
	for _, field := range r.fields {
		fields[field.Name()] = field
	}
	return fields
}

func (r NestedField) PlanModifiers() []string {
	return MapSlice(r.planModifiers, Sentencefy)
}

func (r NestedField) Validators() []string {
	return MapSlice(r.validators, Sentencefy)
}

type Field struct {
	parents  []string
	name     string
	dataType DataType

	required bool
	optional bool
	computed bool

	sensitive bool

	description string
	deprecation string

	writeOnly bool

	defaultDesc   string
	planModifiers []string
	validators    []string
}

func (r Field) Parents() []string {
	return r.parents
}

func (r Field) Name() string {
	return r.name
}

func (r Field) DataType() DataType {
	return r.dataType
}

func (r Field) Required() bool {
	return r.required
}

func (r Field) Optional() bool {
	return r.optional
}

func (r Field) Computed() bool {
	return r.computed
}

func (r Field) Sensitive() bool {
	return r.sensitive
}

func (r Field) Description() string {
	return r.description
}

func (r Field) Deprecation() string {
	return r.deprecation
}

func (r Field) WriteOnly() bool {
	return r.writeOnly
}

func (r Field) Default() string {
	return Sentencefy(r.defaultDesc)
}

func (r Field) PlanModifiers() []string {
	return MapSlice(r.planModifiers, Sentencefy)
}

func (r Field) Validators() []string {
	return MapSlice(r.validators, Sentencefy)
}

func (field Field) Traits() string {
	var traits []string
	traits = append(traits, field.DataType().String())
	if field.Sensitive() {
		traits = append(traits, "Sensitive")
	}
	if field.WriteOnly() {
		traits = append(traits, "[Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)")
	}
	return strings.Join(traits, ", ")
}

func (field Field) nestedKey() string {
	return strings.Join(slices.Concat(field.Parents(), []string{field.Name()}), ".")
}

func (field Field) NestedLink() string {
	switch field.DataType() {
	case DTSingleNestedAttr,
		DTListNestedAttr,
		DTMapNestedAttr,
		DTSetNestedAttr,
		DTObjectAttr,
		DTSingleNestedBlock,
		DTListNestedBlock,
		DTSetNestedBlock:
		return fmt.Sprintf("See the nested schema [here](#nested--%s).", field.nestedKey())
	default:
		return ""
	}
}
