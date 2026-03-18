package metadata

import (
	"bytes"
	"fmt"
	"io"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

type ImportId struct {
	Format  string
	Example string
}

type resourceRenderBuilder struct {
	ProviderName string
	ResourceType string

	Metadata ResourceMetadata

	Subcategory        string
	Examples           []Example
	ObjectDescriptions ObjectDescription

	// Import
	ImportId         *ImportId
	IdentityExamples []Example
}

func (b resourceRenderBuilder) Category() Category {
	return CategoryResource
}

func (b resourceRenderBuilder) renderHeader(w io.Writer) error {
	return renderHeader(w, b.Category(), b.ProviderName, b.ResourceType, b.Subcategory, b.Metadata.Schema.Description)
}

func (b resourceRenderBuilder) renderDescription(w io.Writer) error {
	return renderDescription(w, b.Category(), b.ProviderName, b.ResourceType, b.Metadata.Schema.Deprecation, b.Metadata.Schema.Description)
}

func (b resourceRenderBuilder) renderExample(w io.Writer) error {
	return renderExamples(w, b.Examples)
}

func (b resourceRenderBuilder) renderSchema(w io.Writer) error {
	return renderSchema(w, b.Metadata.Schema.Fields, b.Metadata.Schema.Nested, b.ObjectDescriptions)
}

func (b resourceRenderBuilder) renderImport(w io.Writer) error {
	if identity := b.Metadata.Identity; b.ImportId != nil || identity != nil {
		io.WriteString(w, "## Import\n")

		if b.ImportId != nil {
			io.WriteString(w, "\n")
			if err := b.renderImportId(w, *b.ImportId); err != nil {
				return err
			}
		}

		if identity != nil {
			io.WriteString(w, "\n")
			if err := b.renderImportIdentity(w, *identity); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b resourceRenderBuilder) renderImportId(w io.Writer, importId ImportId) error {
	if _, err := fmt.Fprintf(w, `### Import ID

The [%[1]sterraform import%[1]s command](https://developer.hashicorp.com/terraform/cli/commands/import) can be used with the id format: %[1]s%[2]s%[1]s, for example:

%[1]s%[1]s%[1]sshell
$ terraform import %[3]s.example "%[4]s"
%[1]s%[1]s%[1]s

In Terraform v1.5.0 and later, the [%[1]simport%[1]s block](https://developer.hashicorp.com/terraform/language/block/import) can be used with the %[1]sid%[1]s attribute, for example:

%[1]s%[1]s%[1]sterraform
import {
  to = %[3]s.example
  id = "%[4]s"
}
%[1]s%[1]s%[1]s
`, "`", importId.Format, b.ResourceType, importId.Example); err != nil {
		return err
	}

	return nil
}

func (b resourceRenderBuilder) renderImportIdentity(w io.Writer, schema ResourceIdentitySchema) error {
	formatExample := func(example []byte) []byte {
		return hclwrite.Format(fmt.Appendf(nil, `import {
  to = %s.example
  identity = {
    %s
  }
}`, b.ResourceType, bytes.TrimSpace(example)))
	}

	if _, err := fmt.Fprintf(w, `### Import Identity

In Terraform v1.12.0 and later, the [%[1]simport%[1]s block](https://developer.hashicorp.com/terraform/language/block/import) can be used with the %[1]sidentity%[1]s attribute.
`, "`"); err != nil {
		return err
	}

	for _, example := range b.IdentityExamples {
		if example.Header != nil {
			if _, err := fmt.Fprintf(w, "\n#### Example: %s\n", *example.Header); err != nil {
				return err
			}
		}
		if example.Description != nil {
			if _, err := fmt.Fprintf(w, "\n%s\n", *example.Description); err != nil {
				return err
			}
		}
		if example.HCL != nil {
			if _, err := fmt.Fprintf(w, "\n```terraform\n%s\n```\n", formatExample(example.HCL)); err != nil {
				return err
			}
		}
	}

	sections := []struct {
		name   string
		fields []ResourceIdentityField
	}{
		{
			name:   "Required",
			fields: schema.Fields.RequiredFields(),
		},
		{
			name:   "Optional",
			fields: schema.Fields.OptionalFields(),
		},
	}

	for _, section := range sections {
		if len(section.fields) == 0 {
			continue
		}
		if _, err := fmt.Fprintf(w, `
%s:

`, section.name); err != nil {
			return err
		}

		for _, field := range section.fields {
			if err := b.renderIdentityField(w, field); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b resourceRenderBuilder) renderIdentityField(w io.Writer, field ResourceIdentityField) error {
	if _, err := fmt.Fprintf(w, "- `%s` (%s) %s\n", field.Name, field.Traits(), field.Description); err != nil {
		return err
	}
	return nil
}
