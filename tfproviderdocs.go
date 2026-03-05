package tfproviderdocs

import (
	"context"
	"fmt"
	"io"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/magodo/tfproviderdocs/internal/metadata"
	"github.com/magodo/tfproviderdocs/internal/render"
)

type Generator struct {
	metadata metadata.Metadata
}

func NewGenerator(ctx context.Context, p provider.Provider) (*Generator, error) {
	metadata, diags := metadata.GetMetadata(ctx, p)
	if diags.HasError() {
		return nil, diagsToError(diags)
	}

	return &Generator{metadata: metadata}, nil
}

type Example = render.Example

type ResourceRenderOption struct {
	SubCategory string
	Examples    []Example
}

func (gen Generator) RenderResource(ctx context.Context, w io.Writer, resourceType string, option *ResourceRenderOption) error {
	res, ok := gen.metadata.Resources[resourceType]
	if !ok {
		return fmt.Errorf("Resource type %q not found", resourceType)
	}
	rr := render.ResourceRender{
		ProviderName: gen.metadata.ProviderName,
		ResourceType: resourceType,
		Schema:       res,
	}

	if option != nil {
		rr.Subcategory = option.SubCategory
		rr.Examples = option.Examples
	}
	return rr.Render(w)
}
