package tfproviderdocs

import (
	"context"
	"io"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/magodo/tfproviderdocs/internal/metadata"
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

type Example = metadata.Example
type ImportId = metadata.ImportId
type ResourceRenderOption = metadata.ResourceRenderOption

func (gen Generator) RenderResource(ctx context.Context, w io.Writer, resourceType string, option *ResourceRenderOption) error {
	rr, err := gen.metadata.NewResourceRender(resourceType, option)
	if err != nil {
		return err
	}
	return rr.Execute(w)
}
