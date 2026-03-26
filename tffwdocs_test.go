package tffwdocs_test

import (
	"bytes"
	"context"
	"log"
	"os"
	"testing"
	"text/template"

	tffwdocs "github.com/magodo/terraform-plugin-framework-docs"
	"github.com/magodo/terraform-plugin-framework-docs/internal/metadata"
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
				Header:      "Basic",
				Description: "The basic configuration.",
				HCL: `
provider "examplecloud" {
	name = "foo"
}
	`,
			},
			{
				Header:      "Complete",
				Description: "The complete configuration.",
				HCL: `
provider "examplecloud" {
	name = "foo"
	address = "bar"
	age = 123
	role = "Software Engineer"
}
`,
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
				Header:      "Basic",
				Description: "The basic configuration.",
				HCL: `
resource "examplecloud_resource" "example" {
	name = "foo"
}
`,
			},
			{
				Header:      "Complete",
				Description: "The complete configuration.",
				HCL: `
resource "examplecloud_resource" "example" {
	name = "foo"
	address = "bar"
	age = 123
	role = "Software Engineer"
}
`,
			},
		},
		ImportId: &tffwdocs.ImportId{
			Format:    "<parent_id>/<id>[/<version>]",
			ExampleId: "123/456",
		},
		IdentityExamples: []tffwdocs.Example{
			{
				Header:      "Without Version",
				Description: "Import without version.",
				HCL: `
parent_id = "123"
id = "456"
`,
			},
			{
				Header:      "With Version",
				Description: "Import with version.",
				HCL: `
parent_id = "123"
id = "456"
version = "v2"
`,
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
				Header:      "Basic",
				Description: "The basic configuration.",
				HCL: `
data "examplecloud_resource" "example" {
	name = "foo"
}
	`,
			},
			{
				Header:      "Complete",
				Description: "The complete configuration.",
				HCL: `
data "examplecloud_resource" "example" {
	name = "foo"
	address = "bar"
	age = 123
	role = "Software Engineer"
}
	`,
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
				Header:      "Basic",
				Description: "The basic configuration.",
				HCL: `
ephemeral "examplecloud_resource" "example" {
	name = "foo"
}
	`,
			},
			{
				Header:      "Complete",
				Description: "The complete configuration.",
				HCL: `
ephemeral "examplecloud_resource" "example" {
	name = "foo"
	address = "bar"
	age = 123
	role = "Software Engineer"
}
	`,
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
				Header:      "Basic",
				Description: "The basic configuration.",
				HCL: `
action "examplecloud_resource" "example" {
	config {
		name = "foo"
	}
}
	`,
			},
			{
				Header:      "Complete",
				Description: "The complete configuration.",
				HCL: `
action "examplecloud_resource" "example" {
	config {
		name = "foo"
		address = "bar"
		age = 123
		role = "Software Engineer"
	}
}
	`,
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
				Header:      "Basic",
				Description: "The basic configuration.",
				HCL: `
list "examplecloud_resource" "example" {
	config {
		name = "foo"
	}
}
	`,
			},
			{
				Header:      "Complete",
				Description: "The complete configuration.",
				HCL: `
list "examplecloud_resource" "example" {
	config {
		name = "foo"
		address = "bar"
		age = 123
		role = "Software Engineer"
	}
}
	`,
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
				Header:      "Basic",
				Description: "The basic call.",
				HCL:         `example_function_simple(...)`,
			},
			{
				Header:      "Complete",
				Description: "The complete call.",
				HCL:         `example_function_simple(...)`,
			},
		},
		ReturnDescription: new("This function returns a boolean indicating something."),
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
				Header:      "Basic",
				Description: "The basic call.",
				HCL:         `example_function_retobj(...)`,
			},
			{
				Header:      "Complete",
				Description: "The complete call.",
				HCL:         `example_function_retobj(...)`,
			},
		},
		ReturnDescription: new("This function returns an object indicating something."),
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

func ExampleGenerator_WriteAll() {
	ctx := context.Background()
	gen, err := tffwdocs.NewGenerator(ctx, &testprovider.ExampleCloudProvider{})
	if err != nil {
		log.Fatal(err)
	}
	if err := gen.WriteAll(ctx, "./internal/testprovider/docs", &tffwdocs.RenderOptions{
		Provider: &metadata.ProviderRenderOption{
			Examples: []tffwdocs.Example{
				{
					Header:      "Basic",
					Description: "The basic configuration.",
					HCL: `
provider "examplecloud" {
	name = "foo"
}
	`,
				},
				{
					Header:      "Complete",
					Description: "The complete configuration.",
					HCL: `
provider "examplecloud" {
	name = "foo"
	address = "bar"
	age = 123
	role = "Software Engineer"
}
	`,
				},
			},
		},
		Resources: map[string]tffwdocs.ResourceRenderOption{
			"examplecloud_resource": metadata.ResourceRenderOption{
				Subcategory: "abc",
				Examples: []tffwdocs.Example{
					{
						Header:      "Basic",
						Description: "The basic configuration.",
						HCL: `
resource "examplecloud_resource" "example" {
	name = "foo"
}
`,
					},
					{
						Header:      "Complete",
						Description: "The complete configuration.",
						HCL: `
resource "examplecloud_resource" "example" {
	name = "foo"
	address = "bar"
	age = 123
	role = "Software Engineer"
}
`,
					},
				},
				ImportId: &tffwdocs.ImportId{
					Format:    "<parent_id>/<id>[/<version>]",
					ExampleId: "123/456",
					ExampleBlk: `
import {
	to = examplecloud_resource.example
	id = "123/456"
}
`,
				},
				IdentityExamples: []tffwdocs.Example{
					{
						Header:      "Without Version",
						Description: "Import without version.",
						HCL: `
import {
	to = examplecloud_resource.example
	identity = {
		parent_id = "123"
		id = "456"
	}
}
`,
					},
					{
						Header:      "With Version",
						Description: "Import with version.",
						HCL: `
import {
	to = examplecloud_resource.example
	identity = {
		parent_id = "123"
		id = "456"
		version = "v2"
	}
}
`,
					},
				},
			},
		},
	}); err != nil {
		log.Fatal(err)
	}
	// Output:
}
