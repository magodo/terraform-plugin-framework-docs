package tfproviderdocs

import (
	"bytes"
	"os"
	"testing"

	"github.com/magodo/tfproviderdocs/internal/testprovider"
	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	g, err := NewGenerator(t.Context(), &testprovider.ExampleCloudProvider{})
	require.NoError(t, err)
	var buf bytes.Buffer
	require.NoError(t, g.RenderResource(t.Context(), &buf, "examplecloud_resource", &ResourceRenderOption{
		SubCategory: "abc",
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
	}))
	expected, err := os.ReadFile("./testdata/resource.md")
	require.NoError(t, err)
	require.Equal(t, string(expected), buf.String())
}
