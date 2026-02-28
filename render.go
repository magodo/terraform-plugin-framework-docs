package tfproviderdocs

import (
	"fmt"
	"io"
	"text/template"
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

const TplProperty = "- `{{ .Name }}` ({{ .Traits }}) {{ .Description }}" +
	`{{ if .Default }} {{ .Default }}{{ end }}` +
	`{{ if .NestedLink }} {{ .NestedLink }}{{ end }}` +
	`{{- template "planmodifier-indent" . }}` +
	`{{- template "validator-indent" . }}` +
	`{{- template "deprecation-indent" . }}`

const TplProperties = `{{ range . }}
{{ template "property" . -}}
{{ end }}`

const TplDeprecationIndent = `{{ with .Deprecation }}

	!> {{ . }}
{{ end }}`

const TplPlanModifier = `{{ with .PlanModifiers }}

Plan Modifiers:
{{ range . }}
- {{ . -}}
{{ end -}}
{{ end }}`

const TplPlanModifierIndent = `{{ with .PlanModifiers }}

	Plan Modifiers:
{{ range . }}
	- {{ . -}}
{{ end -}}
{{ end }}`

const TplValidator = `{{ with .Validators }}

Validators:
{{ range . }}
- {{ . -}}
{{ end -}}
{{ end }}`

const TplValidatorIndent = `{{ with .Validators }}

	Validators:
{{ range . }}
	- {{ . -}}
{{ end -}}
{{ end }}`

var TplNested = fmt.Sprintf(`{{- range $key, $value := . }}
<a id="nested--{{ $key }}"></a>
### Nested Schema for %[1]s{{ $key }}%[1]s

{{- template "planmodifier" . }} {{- template "validator" . }}

{{- with $value.Infos.RequiredInfos }}

Required:
{{ template "properties" . -}}
{{ end }}

{{- with $value.Infos.OptionalInfos }}

Optional:
{{ template "properties" . -}}
{{ end }}

{{- with $value.Infos.ComputedInfos }}

Read-Only:
{{ template "properties" . -}}
{{ end }}
{{ end }})`, "`")

const TplSchema = `## Schema

{{- with .Infos.RequiredInfos }}

### Required
{{ template "properties" . -}}
{{ end }}

{{- with .Infos.OptionalInfos }}

### Optional
{{ template "properties" . -}}
{{ end }}

{{- with .Infos.ComputedInfos }}

### Read-Only
{{ template "properties" . -}}
{{ end }}

{{- with .Nested }}
{{ template "nested" . }}
{{ end }}
`

func (render ResourceRender) Execute(w io.Writer) error {
	tpl := template.Must(template.New("schema").Parse(TplSchema))
	template.Must(tpl.New("properties").Parse(TplProperties))
	template.Must(tpl.New("property").Parse(TplProperty))
	template.Must(tpl.New("deprecation-indent").Parse(TplDeprecationIndent))
	template.Must(tpl.New("planmodifier").Parse(TplPlanModifier))
	template.Must(tpl.New("planmodifier-indent").Parse(TplPlanModifierIndent))
	template.Must(tpl.New("validator").Parse(TplValidator))
	template.Must(tpl.New("validator-indent").Parse(TplValidatorIndent))
	template.Must(tpl.New("nested").Parse(TplNested))
	if err := tpl.Execute(w, render.SchemaInfo); err != nil {
		return err
	}
	return nil
}
