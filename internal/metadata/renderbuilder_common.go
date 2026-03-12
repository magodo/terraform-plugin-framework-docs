package metadata

import (
	"bytes"
	"fmt"
	"io"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

type Example struct {
	Header      *string
	Description *string
	HCL         []byte
}

func renderExamples(w io.Writer, examples []Example) error {
	if len(examples) != 0 {
		if _, err := io.WriteString(w, "## Example Usage\n"); err != nil {
			return err
		}
		for _, example := range examples {
			if example.Header != nil {
				if _, err := fmt.Fprintf(w, "\n### %s\n", *example.Header); err != nil {
					return err
				}
			}
			if example.Description != nil {
				if _, err := fmt.Fprintf(w, "\n%s\n", *example.Description); err != nil {
					return err
				}
			}
			if example.HCL != nil {
				if _, err := fmt.Fprintf(w, "\n```terraform\n%s\n```\n", bytes.TrimSpace(hclwrite.Format(example.HCL))); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
