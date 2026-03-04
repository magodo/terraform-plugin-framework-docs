package schema

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

type Fields map[string]Field

func (infos Fields) RequiredFields() []Field {
	var out []Field
	for _, info := range infos {
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

func (infos Fields) OptionalFields() []Field {
	var out []Field
	for _, info := range infos {
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

func (infos Fields) ComputedFields() []Field {
	var out []Field
	for _, info := range infos {
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

	PlanModifiers []string
	Validators    []string
	DefaultDesc   *string

	WriteOnly bool
}

func (info Field) NestedKey() string {
	return strings.Join(slices.Concat(info.Parents, []string{info.Name}), ".")
}

func (info Field) NestedLink() string {
	switch info.DataType {
	case DTSingleNestedAttr,
		DTListNestedAttr,
		DTMapNestedAttr,
		DTSetNestedAttr,
		DTObjectAttr,
		DTSingleNestedBlock,
		DTListNestedBlock,
		DTSetNestedBlock:
		return fmt.Sprintf("See the nested schema [here](#nested--%s).", info.NestedKey())
	default:
		return ""
	}
}

func (info Field) Default() string {
	if info.DefaultDesc == nil {
		return ""
	}
	return Sentencefy(*info.DefaultDesc)

}

func (info Field) Traits() string {
	var traits []string
	traits = append(traits, info.DataType.String())
	if info.Sensitive {
		traits = append(traits, "Sensitive")
	}
	if info.WriteOnly {
		traits = append(traits, "[Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)")
	}
	return strings.Join(traits, ", ")
}
