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

func (infos SchemaInfos) RequiredInfos() []SchemaInfo {
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

func (infos SchemaInfos) OptionalInfos() []SchemaInfo {
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

func (infos SchemaInfos) ComputedInfos() []SchemaInfo {
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
	PlanModifiers []string
	Validators    []string
	Infos         SchemaInfos
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

	PlanModifiers []string
	Validators    []string
	DefaultDesc   *string

	WriteOnly bool
}
