package tfproviderdocs

import (
	"html/template"
	"io"
)

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

const PropertyTemplate = "- `{{ .Name }}` ({{ .DataType }}) {{ .Description }}"
const PropertiesTemplate = `{{ range . }}
{{ template "property" . -}}
{{ end }}`

const SchemaTemplate = `## Schema

{{ with .Infos.Requireds -}}
### Required
{{ template "properties" . -}}
{{ end }}

{{ with .Infos.Optionals -}}
### Optional
{{ template "properties" . -}}
{{ end }}

{{ with .Infos.Computeds -}}
### Computed
{{ template "properties" . -}}
{{ end }}
`

func (render ResourceRender) Execute(w io.Writer) error {
	tpl := template.Must(template.New("schema").Parse(SchemaTemplate))
	template.Must(tpl.New("properties").Parse(PropertiesTemplate))
	template.Must(tpl.New("property").Parse(PropertyTemplate))
	if err := tpl.Execute(w, render.SchemaInfo); err != nil {
		return err
	}
	return nil
}
