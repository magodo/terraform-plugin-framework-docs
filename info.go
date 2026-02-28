package tfproviderdocs

import (
	"sort"
)

type ResourceInfos map[string]ResourceInfo

type ResourceInfo struct {
	Description string
	Deprecation string

	Infos SchemaInfos

	// Including nested attribute object or block object.
	Nested NestedSchemaInfos
}

type SchemaInfos map[string]SchemaInfo

func (infos SchemaInfos) Requireds() []SchemaInfo {
	var out []SchemaInfo
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

func (infos SchemaInfos) Optionals() []SchemaInfo {
	var out []SchemaInfo
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

func (infos SchemaInfos) Computeds() []SchemaInfo {
	var out []SchemaInfo
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

type NestedSchemaInfos map[string]NestedSchemaInfo

type NestedSchemaInfo struct {
	PlanModifierDescriptions []string
	Infos                    SchemaInfos
}

type SchemaInfo struct {
	Parents  []string
	Name     string
	DataType DataType

	Required bool
	Optional bool
	Computed bool

	Sensitive bool

	Description string
	Deprecation string

	PlanModifierDescriptions []string
	DefaultDesc              *string

	WriteOnly bool
}
