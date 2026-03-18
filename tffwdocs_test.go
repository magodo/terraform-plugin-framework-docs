package tffwdocs_test

import (
	"bytes"
	"os"
	"testing"
	"text/template"

	tffwdocs "github.com/magodo/terraform-plugin-framework-docs"
	"github.com/magodo/terraform-plugin-framework-docs/internal/testprovider"
	"github.com/stretchr/testify/require"
)

func TestProviderRender(t *testing.T) {
	g, err := tffwdocs.NewGenerator(t.Context(), &testprovider.ExampleCloudProvider{})
	require.NoError(t, err)

	// Render the minimal version
	var buf bytes.Buffer
	require.NoError(t, g.RenderProvider(t.Context(), &buf, nil))
	expected, err := os.ReadFile("./testdata/provider_minimal.md")
	require.NoError(t, err)
	require.Equal(t, string(expected), buf.String(), "minimal")

	// Render the complete version
	opt := &tffwdocs.ProviderRenderOption{
		Examples: []tffwdocs.Example{
			{
				Header:      new("Basic"),
				Description: new("The basic configuration."),
				HCL: []byte(`
provider "examplecloud_resource" "example" {
	name = "foo"
}
	`),
			},
			{
				Header:      new("Complete"),
				Description: new("The complete configuration."),
				HCL: []byte(`
provider "examplecloud_resource" "example" {
	name = "foo"
	address = "bar"
	age = 123
	role = "Software Engineer"
}
	`),
			},
		},
	}
	buf = *bytes.NewBuffer(nil)
	require.NoError(t, g.RenderProvider(t.Context(), &buf, opt))
	expected, err = os.ReadFile("./testdata/provider_complete.md")
	require.NoError(t, err)
	require.Equal(t, string(expected), buf.String(), "complete")

	// Render the custom template
	opt.Template = template.Must(template.New("").Parse(`{{ .Header }}
{{ .Description }}
Some note...

{{ .Example }}`))
	buf = *bytes.NewBuffer(nil)
	require.NoError(t, g.RenderProvider(t.Context(), &buf, opt))
	expected, err = os.ReadFile("./testdata/provider_custom.md")
	require.NoError(t, err)
	require.Equal(t, string(expected), buf.String(), "custom")
}

func TestResourceRender(t *testing.T) {
	g, err := tffwdocs.NewGenerator(t.Context(), &testprovider.ExampleCloudProvider{})
	require.NoError(t, err)

	// Render the minimal version
	var buf bytes.Buffer
	require.NoError(t, g.RenderResource(t.Context(), &buf, "examplecloud_resource", nil))
	expected, err := os.ReadFile("./testdata/resource_minimal.md")
	require.NoError(t, err)
	require.Equal(t, string(expected), buf.String(), "minimal")

	// Render the complete version
	opt := &tffwdocs.ResourceRenderOption{
		Subcategory: "abc",
		Examples: []tffwdocs.Example{
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
		ImportId: &tffwdocs.ImportId{
			Format:  "<parent_id>/<id>[/<version>]",
			Example: "123/456",
		},
		IdentityExamples: []tffwdocs.Example{
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
		ObjectDescription: tffwdocs.ObjectDescription{
			"object": map[string]string{
				"foo": "A foo field.",
			},
			"object.foo": map[string]string{
				"bar": "A bar field.",
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

func TestDataSourceRender(t *testing.T) {
	g, err := tffwdocs.NewGenerator(t.Context(), &testprovider.ExampleCloudProvider{})
	require.NoError(t, err)

	// Render the minimal version
	var buf bytes.Buffer
	require.NoError(t, g.RenderDataSource(t.Context(), &buf, "examplecloud_resource", nil))
	expected, err := os.ReadFile("./testdata/datasource_minimal.md")
	require.NoError(t, err)
	require.Equal(t, string(expected), buf.String(), "minimal")

	// Render the complete version
	opt := &tffwdocs.DataSourceRenderOption{
		Subcategory: "abc",
		Examples: []tffwdocs.Example{
			{
				Header:      new("Basic"),
				Description: new("The basic configuration."),
				HCL: []byte(`
data "examplecloud_resource" "example" {
	name = "foo"
}
	`),
			},
			{
				Header:      new("Complete"),
				Description: new("The complete configuration."),
				HCL: []byte(`
data "examplecloud_resource" "example" {
	name = "foo"
	address = "bar"
	age = 123
	role = "Software Engineer"
}
	`),
			},
		},
	}
	buf = *bytes.NewBuffer(nil)
	require.NoError(t, g.RenderDataSource(t.Context(), &buf, "examplecloud_resource", opt))
	expected, err = os.ReadFile("./testdata/datasource_complete.md")
	require.NoError(t, err)
	require.Equal(t, string(expected), buf.String(), "complete")

	// Render the custom template
	opt.Template = template.Must(template.New("").Parse(`{{ .Header }}
{{ .Description }}
Some note...

{{ .Example }}`))
	buf = *bytes.NewBuffer(nil)
	require.NoError(t, g.RenderDataSource(t.Context(), &buf, "examplecloud_resource", opt))
	expected, err = os.ReadFile("./testdata/datasource_custom.md")
	require.NoError(t, err)
	require.Equal(t, string(expected), buf.String(), "custom")
}

func TestEphemeralRender(t *testing.T) {
	g, err := tffwdocs.NewGenerator(t.Context(), &testprovider.ExampleCloudProvider{})
	require.NoError(t, err)

	// Render the minimal version
	var buf bytes.Buffer
	require.NoError(t, g.RenderEphemeralResource(t.Context(), &buf, "examplecloud_resource", nil))
	expected, err := os.ReadFile("./testdata/ephemeral_minimal.md")
	require.NoError(t, err)
	require.Equal(t, string(expected), buf.String(), "minimal")

	// Render the complete version
	opt := &tffwdocs.EphemeralResourceRenderOption{
		Subcategory: "abc",
		Examples: []tffwdocs.Example{
			{
				Header:      new("Basic"),
				Description: new("The basic configuration."),
				HCL: []byte(`
ephemeral "examplecloud_resource" "example" {
	name = "foo"
}
	`),
			},
			{
				Header:      new("Complete"),
				Description: new("The complete configuration."),
				HCL: []byte(`
ephemeral "examplecloud_resource" "example" {
	name = "foo"
	address = "bar"
	age = 123
	role = "Software Engineer"
}
	`),
			},
		},
	}
	buf = *bytes.NewBuffer(nil)
	require.NoError(t, g.RenderEphemeralResource(t.Context(), &buf, "examplecloud_resource", opt))
	expected, err = os.ReadFile("./testdata/ephemeral_complete.md")
	require.NoError(t, err)
	require.Equal(t, string(expected), buf.String(), "complete")

	// Render the custom template
	opt.Template = template.Must(template.New("").Parse(`{{ .Header }}
{{ .Description }}
Some note...

{{ .Example }}`))
	buf = *bytes.NewBuffer(nil)
	require.NoError(t, g.RenderEphemeralResource(t.Context(), &buf, "examplecloud_resource", opt))
	expected, err = os.ReadFile("./testdata/ephemeral_custom.md")
	require.NoError(t, err)
	require.Equal(t, string(expected), buf.String(), "custom")
}

func TestActionRender(t *testing.T) {
	g, err := tffwdocs.NewGenerator(t.Context(), &testprovider.ExampleCloudProvider{})
	require.NoError(t, err)

	// Render the minimal version
	var buf bytes.Buffer
	require.NoError(t, g.RenderAction(t.Context(), &buf, "examplecloud_resource", nil))
	expected, err := os.ReadFile("./testdata/action_minimal.md")
	require.NoError(t, err)
	require.Equal(t, string(expected), buf.String(), "minimal")

	// Render the complete version
	opt := &tffwdocs.ActionRenderOption{
		Subcategory: "abc",
		Examples: []tffwdocs.Example{
			{
				Header:      new("Basic"),
				Description: new("The basic configuration."),
				HCL: []byte(`
action "examplecloud_resource" "example" {
	config {
		name = "foo"
	}
}
	`),
			},
			{
				Header:      new("Complete"),
				Description: new("The complete configuration."),
				HCL: []byte(`
action "examplecloud_resource" "example" {
	config {
		name = "foo"
		address = "bar"
		age = 123
		role = "Software Engineer"
	}
}
	`),
			},
		},
	}
	buf = *bytes.NewBuffer(nil)
	require.NoError(t, g.RenderAction(t.Context(), &buf, "examplecloud_resource", opt))
	expected, err = os.ReadFile("./testdata/action_complete.md")
	require.NoError(t, err)
	require.Equal(t, string(expected), buf.String(), "complete")

	// Render the custom template
	opt.Template = template.Must(template.New("").Parse(`{{ .Header }}
{{ .Description }}
Some note...

{{ .Example }}`))
	buf = *bytes.NewBuffer(nil)
	require.NoError(t, g.RenderAction(t.Context(), &buf, "examplecloud_resource", opt))
	expected, err = os.ReadFile("./testdata/action_custom.md")
	require.NoError(t, err)
	require.Equal(t, string(expected), buf.String(), "custom")
}

func TestListRender(t *testing.T) {
	g, err := tffwdocs.NewGenerator(t.Context(), &testprovider.ExampleCloudProvider{})
	require.NoError(t, err)

	// Render the minimal version
	var buf bytes.Buffer
	require.NoError(t, g.RenderListResource(t.Context(), &buf, "examplecloud_resource", nil))
	expected, err := os.ReadFile("./testdata/list_minimal.md")
	require.NoError(t, err)
	require.Equal(t, string(expected), buf.String(), "minimal")

	// Render the complete version
	opt := &tffwdocs.ListResourceRenderOption{
		Subcategory: "abc",
		Examples: []tffwdocs.Example{
			{
				Header:      new("Basic"),
				Description: new("The basic configuration."),
				HCL: []byte(`
list "examplecloud_resource" "example" {
	config {
		name = "foo"
	}
}
	`),
			},
			{
				Header:      new("Complete"),
				Description: new("The complete configuration."),
				HCL: []byte(`
list "examplecloud_resource" "example" {
	config {
		name = "foo"
		address = "bar"
		age = 123
		role = "Software Engineer"
	}
}
	`),
			},
		},
	}
	buf = *bytes.NewBuffer(nil)
	require.NoError(t, g.RenderListResource(t.Context(), &buf, "examplecloud_resource", opt))
	expected, err = os.ReadFile("./testdata/list_complete.md")
	require.NoError(t, err)
	require.Equal(t, string(expected), buf.String(), "complete")

	// Render the custom template
	opt.Template = template.Must(template.New("").Parse(`{{ .Header }}
{{ .Description }}
Some note...

{{ .Example }}`))
	buf = *bytes.NewBuffer(nil)
	require.NoError(t, g.RenderListResource(t.Context(), &buf, "examplecloud_resource", opt))
	expected, err = os.ReadFile("./testdata/list_custom.md")
	require.NoError(t, err)
	require.Equal(t, string(expected), buf.String(), "custom")
}

func TestFunctionRenderSimple(t *testing.T) {
	g, err := tffwdocs.NewGenerator(t.Context(), &testprovider.ExampleCloudProvider{})
	require.NoError(t, err)

	// Render the minimal version
	var buf bytes.Buffer
	require.NoError(t, g.RenderFunction(t.Context(), &buf, "example_function_simple", nil))
	expected, err := os.ReadFile("./testdata/function_simple_minimal.md")
	require.NoError(t, err)
	require.Equal(t, string(expected), buf.String(), "minimal")

	// Render the complete version
	opt := &tffwdocs.FunctionRenderOption{
		Subcategory: "abc",
		Examples: []tffwdocs.Example{
			{
				Header:      new("Basic"),
				Description: new("The basic call."),
				HCL:         []byte(`example_function_simple(...)`),
			},
			{
				Header:      new("Complete"),
				Description: new("The complete call."),
				HCL:         []byte(`example_function_simple(...)`),
			},
		},
		ReturnDescription: new("This function returns a boolean indicating something."),
		ObjectDescription: tffwdocs.ObjectDescription{
			"object": map[string]string{
				"foo": "The foo field.",
				"bar": "The bar field.",
			},
		},
	}
	buf = *bytes.NewBuffer(nil)
	require.NoError(t, g.RenderFunction(t.Context(), &buf, "example_function_simple", opt))
	expected, err = os.ReadFile("./testdata/function_simple_complete.md")
	require.NoError(t, err)
	require.Equal(t, string(expected), buf.String(), "complete")

	// Render the custom template
	opt.Template = template.Must(template.New("").Parse(`{{ .Header }}
{{ .Description }}
Some note...
{{ .Example }}`))
	buf = *bytes.NewBuffer(nil)
	require.NoError(t, g.RenderFunction(t.Context(), &buf, "example_function_simple", opt))
	expected, err = os.ReadFile("./testdata/function_simple_custom.md")
	require.NoError(t, err)
	require.Equal(t, string(expected), buf.String(), "custom")
}

func TestFunctionRenderRetObj(t *testing.T) {
	g, err := tffwdocs.NewGenerator(t.Context(), &testprovider.ExampleCloudProvider{})
	require.NoError(t, err)

	// Render the minimal version
	var buf bytes.Buffer
	require.NoError(t, g.RenderFunction(t.Context(), &buf, "example_function_retobj", nil))
	expected, err := os.ReadFile("./testdata/function_retobj_minimal.md")
	require.NoError(t, err)
	require.Equal(t, string(expected), buf.String(), "minimal")

	// Render the complete version
	opt := &tffwdocs.FunctionRenderOption{
		Subcategory: "abc",
		Examples: []tffwdocs.Example{
			{
				Header:      new("Basic"),
				Description: new("The basic call."),
				HCL:         []byte(`example_function_retobj(...)`),
			},
			{
				Header:      new("Complete"),
				Description: new("The complete call."),
				HCL:         []byte(`example_function_retobj(...)`),
			},
		},
		ReturnDescription: new("This function returns an object indicating something."),
		ObjectDescription: tffwdocs.ObjectDescription{
			"object": map[string]string{
				"foo": "The foo field.",
			},
			"object.foo": map[string]string{
				"bar": "The bar field.",
			},
			"": map[string]string{
				"retfoo": "The retfoo field.",
			},
			"retfoo": map[string]string{
				"retbar": "The retbar field.",
			},
		},
	}
	buf = *bytes.NewBuffer(nil)
	require.NoError(t, g.RenderFunction(t.Context(), &buf, "example_function_retobj", opt))
	expected, err = os.ReadFile("./testdata/function_retobj_complete.md")
	require.NoError(t, err)
	require.Equal(t, string(expected), buf.String(), "complete")

	// Render the custom template
	opt.Template = template.Must(template.New("").Parse(`{{ .Header }}
{{ .Description }}
Some note...
{{ .Example }}`))
	buf = *bytes.NewBuffer(nil)
	require.NoError(t, g.RenderFunction(t.Context(), &buf, "example_function_retobj", opt))
	expected, err = os.ReadFile("./testdata/function_retobj_custom.md")
	require.NoError(t, err)
	require.Equal(t, string(expected), buf.String(), "custom")
}
