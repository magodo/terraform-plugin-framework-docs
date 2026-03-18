package tffwdocs

import (
	"context"
	"io"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/magodo/terraform-plugin-framework-docs/internal/metadata"
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
type ObjectDescription = metadata.ObjectDescription
type ProviderRenderOption = metadata.ProviderRenderOption
type ResourceRenderOption = metadata.ResourceRenderOption
type DataSourceRenderOption = metadata.DataSourceRenderOption
type EphemeralResourceRenderOption = metadata.EphemeralRenderOption
type ActionRenderOption = metadata.ActionRenderOption
type ListResourceRenderOption = metadata.ListRenderOption
type FunctionRenderOption = metadata.FunctionRenderOption

func (gen Generator) RenderProvider(ctx context.Context, w io.Writer, option *ProviderRenderOption) error {
	rr, err := gen.metadata.NewProviderRender(option)
	if err != nil {
		return err
	}
	return rr.Execute(w)
}

func (gen Generator) RenderResource(ctx context.Context, w io.Writer, resourceType string, option *ResourceRenderOption) error {
	rr, err := gen.metadata.NewResourceRender(resourceType, option)
	if err != nil {
		return err
	}
	return rr.Execute(w)
}

func (gen Generator) RenderDataSource(ctx context.Context, w io.Writer, dataSourceType string, option *DataSourceRenderOption) error {
	rr, err := gen.metadata.NewDataSourceRender(dataSourceType, option)
	if err != nil {
		return err
	}
	return rr.Execute(w)
}

func (gen Generator) RenderEphemeralResource(ctx context.Context, w io.Writer, ephemeralResourceType string, option *EphemeralResourceRenderOption) error {
	rr, err := gen.metadata.NewEphemeralRender(ephemeralResourceType, option)
	if err != nil {
		return err
	}
	return rr.Execute(w)
}

func (gen Generator) RenderAction(ctx context.Context, w io.Writer, actionType string, option *ActionRenderOption) error {
	rr, err := gen.metadata.NewActionRender(actionType, option)
	if err != nil {
		return err
	}
	return rr.Execute(w)
}

func (gen Generator) RenderListResource(ctx context.Context, w io.Writer, listResourceType string, option *ListResourceRenderOption) error {
	rr, err := gen.metadata.NewListRender(listResourceType, option)
	if err != nil {
		return err
	}
	return rr.Execute(w)
}

func (gen Generator) RenderFunction(ctx context.Context, w io.Writer, functionName string, option *FunctionRenderOption) error {
	rr, err := gen.metadata.NewFunctionRender(functionName, option)
	if err != nil {
		return err
	}
	return rr.Execute(w)
}
