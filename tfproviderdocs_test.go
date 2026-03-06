package tfproviderdocs

import (
	"bytes"
	"os"
	"testing"
	"text/template"

	"github.com/magodo/tfproviderdocs/internal/testprovider"
	"github.com/stretchr/testify/require"
)

func TestResourceRender(t *testing.T) {
	g, err := NewGenerator(t.Context(), &testprovider.ExampleCloudProvider{})
	require.NoError(t, err)

	// Render the minimal version
	var buf bytes.Buffer
	require.NoError(t, g.RenderResource(t.Context(), &buf, "examplecloud_resource", nil))
	expected, err := os.ReadFile("./testdata/resource_minimal.md")
	require.NoError(t, err)
	require.Equal(t, string(expected), buf.String(), "minimal")

	// Render the complete version
	opt := &ResourceRenderOption{
		Subcategory: "abc",
		Examples: []Example{
			{
				Header:      new("Basic"),
				Description: new("The basic configuration."),
				HCL: []byte(`
resource "examplecloud_resource" "example" {
	name = "foo"
}
`),
			},
			{
				Header:      new("Complete"),
				Description: new("The complete configuration."),
				HCL: []byte(`
resource "examplecloud_resource" "example" {
	name = "foo"
	address = "bar"
	age = 123
	role = "Software Engineer"
}
`),
			},
		},
		ImportId: &ImportId{
			Format:  "<parent_id>/<id>[/<version>]",
			Example: "123/456",
		},
		IdentityExamples: []Example{
			{
				Header:      new("Without Version"),
				Description: new("Import without version."),
				HCL: []byte(`
parent_id = "123"
id = "456"
`),
			},
			{
				Header:      new("With Version"),
				Description: new("Import with version."),
				HCL: []byte(`
parent_id = "123"
id = "456"
version = "v2"
`),
			},
		},
	}
	buf = *bytes.NewBuffer(nil)
	require.NoError(t, g.RenderResource(t.Context(), &buf, "examplecloud_resource", opt))
	expected, err = os.ReadFile("./testdata/resource_complete.md")
	require.NoError(t, err)
	require.Equal(t, string(expected), buf.String(), "complete")

	// Render the custom template
	opt.Template = template.Must(template.New("").Parse(`{{ .Header }}
{{ .Description }}
Some note...

{{ .Example }}`))
	buf = *bytes.NewBuffer(nil)
	require.NoError(t, g.RenderResource(t.Context(), &buf, "examplecloud_resource", opt))
	expected, err = os.ReadFile("./testdata/resource_custom.md")
	require.NoError(t, err)
	require.Equal(t, string(expected), buf.String(), "custom")
}
