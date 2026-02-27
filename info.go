package tfproviderdocs

type ResourceInfos map[string]ResourceInfo

type ResourceInfo struct {
	Description string
	Deprecation string

	Infos SchemaInfos

	// Including nested attribute object or block object.
	Nested NestedSchemaInfos
}

type SchemaInfos map[string]SchemaInfo
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
