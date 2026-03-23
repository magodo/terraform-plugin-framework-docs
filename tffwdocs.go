package tffwdocs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

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
type ProviderRenderOption = metadata.ProviderRenderOption
type ResourceRenderOption = metadata.ResourceRenderOption
type DataSourceRenderOption = metadata.DataSourceRenderOption
type EphemeralResourceRenderOption = metadata.EphemeralRenderOption
type ActionRenderOption = metadata.ActionRenderOption
type ListResourceRenderOption = metadata.ListRenderOption
type FunctionRenderOption = metadata.FunctionRenderOption

type RenderOptions struct {
	Provider           *ProviderRenderOption
	Resources          map[string]ResourceRenderOption
	DataSources        map[string]DataSourceRenderOption
	EphemeralResources map[string]EphemeralResourceRenderOption
	ListResources      map[string]ListResourceRenderOption
	Actions            map[string]ActionRenderOption
	Functions          map[string]FunctionRenderOption
}

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

// WriteAll generates the documents for the provider and all registered resources, actions, functions, etc.
// and write to the specified document directory, which should be existing.
func (gen Generator) WriteAll(ctx context.Context, docDir string, opts *RenderOptions) error {
	if opts == nil {
		opts = &RenderOptions{
			Provider:           nil,
			Resources:          map[string]ResourceRenderOption{},
			DataSources:        map[string]DataSourceRenderOption{},
			EphemeralResources: map[string]EphemeralResourceRenderOption{},
			ListResources:      map[string]ListResourceRenderOption{},
			Actions:            map[string]ActionRenderOption{},
			Functions:          map[string]FunctionRenderOption{},
		}
	}

	// Provider
	{
		var buf bytes.Buffer
		if err := gen.RenderProvider(ctx, &buf, opts.Provider); err != nil {
			return fmt.Errorf("render provider: %v", err)
		}
		if err := os.WriteFile(filepath.Join(docDir, "index.md"), buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("write provider file: %v", err)
		}
	}

	// Resources
	for name := range gen.metadata.Resources {
		var buf bytes.Buffer
		var opt *ResourceRenderOption
		if optt, ok := opts.Resources[name]; ok {
			opt = &optt
		}
		if err := gen.RenderResource(ctx, &buf, name, opt); err != nil {
			return fmt.Errorf("render resource %q: %v", name, err)
		}
		dir := filepath.Join(docDir, "resources")
		if err := os.Mkdir(dir, 0755); err != nil && !os.IsExist(err) {
			return fmt.Errorf("mkdir %q: %v", dir, err)
		}
		if err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("%s.md", name)), buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("write resource file for %q: %v", name, err)
		}
	}

	// DataSource
	for name := range gen.metadata.DataSources {
		var buf bytes.Buffer
		var opt *DataSourceRenderOption
		if optt, ok := opts.DataSources[name]; ok {
			opt = &optt
		}
		if err := gen.RenderDataSource(ctx, &buf, name, opt); err != nil {
			return fmt.Errorf("render data source %q: %v", name, err)
		}
		dir := filepath.Join(docDir, "data-sources")
		if err := os.Mkdir(dir, 0755); err != nil && !os.IsExist(err) {
			return fmt.Errorf("mkdir %q: %v", dir, err)
		}
		if err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("%s.md", name)), buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("write data source file for %q: %v", name, err)
		}
	}

	// Ephemeral Resource
	for name := range gen.metadata.Ephemerals {
		var buf bytes.Buffer
		var opt *EphemeralResourceRenderOption
		if optt, ok := opts.EphemeralResources[name]; ok {
			opt = &optt
		}
		if err := gen.RenderEphemeralResource(ctx, &buf, name, opt); err != nil {
			return fmt.Errorf("render ephemeral resource %q: %v", name, err)
		}
		dir := filepath.Join(docDir, "ephemeral-resources")
		if err := os.Mkdir(dir, 0755); err != nil && !os.IsExist(err) {
			return fmt.Errorf("mkdir %q: %v", dir, err)
		}
		if err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("%s.md", name)), buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("write ephemeral resource file for %q: %v", name, err)
		}
	}

	// List Resource
	for name := range gen.metadata.Lists {
		var buf bytes.Buffer
		var opt *ListResourceRenderOption
		if optt, ok := opts.ListResources[name]; ok {
			opt = &optt
		}
		if err := gen.RenderListResource(ctx, &buf, name, opt); err != nil {
			return fmt.Errorf("render list resource %q: %v", name, err)
		}
		dir := filepath.Join(docDir, "list-resources")
		if err := os.Mkdir(dir, 0755); err != nil && !os.IsExist(err) {
			return fmt.Errorf("mkdir %q: %v", dir, err)
		}
		if err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("%s.md", name)), buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("write list resource file for %q: %v", name, err)
		}
	}

	// Action
	for name := range gen.metadata.Actions {
		var buf bytes.Buffer
		var opt *ActionRenderOption
		if optt, ok := opts.Actions[name]; ok {
			opt = &optt
		}
		if err := gen.RenderAction(ctx, &buf, name, opt); err != nil {
			return fmt.Errorf("render action %q: %v", name, err)
		}
		dir := filepath.Join(docDir, "actions")
		if err := os.Mkdir(dir, 0755); err != nil && !os.IsExist(err) {
			return fmt.Errorf("mkdir %q: %v", dir, err)
		}
		if err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("%s.md", name)), buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("write action file for %q: %v", name, err)
		}
	}

	// Function
	for name := range gen.metadata.Functions {
		var buf bytes.Buffer
		var opt *FunctionRenderOption
		if optt, ok := opts.Functions[name]; ok {
			opt = &optt
		}
		if err := gen.RenderFunction(ctx, &buf, name, opt); err != nil {
			return fmt.Errorf("render function %q: %v", name, err)
		}
		dir := filepath.Join(docDir, "functions")
		if err := os.Mkdir(dir, 0755); err != nil && !os.IsExist(err) {
			return fmt.Errorf("mkdir %q: %v", dir, err)
		}
		if err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("%s.md", name)), buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("write function file for %q: %v", name, err)
		}
	}

	return nil
}
