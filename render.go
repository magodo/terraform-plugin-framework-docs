package tfproviderdocs

import "io"

type ResourceRender struct {
	ProviderName string
	ResourceType string
	Subcategory  string

	Example    string
	SchemaInfo ResourceInfo

	// TODO
	// IdentitySchemaInfo string
	// ImportId           string
}

func (render ResourceRender) Execute(w io.Writer) {
}
