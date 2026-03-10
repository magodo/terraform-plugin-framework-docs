package metadata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type Metadata struct {
	ProviderName string
	Resources    ResourceMetadatas
	DataSources  DataSourceMetadatas
}

type ResourceMetadatas map[string]ResourceMetadata

type ResourceMetadata struct {
	Schema   ResourceSchema
	Identity *ResourceIdentitySchema
}

type DataSourceMetadatas map[string]DataSourceMetadata

type DataSourceMetadata struct {
	Schema DataSourceSchema
}

func GetMetadata(ctx context.Context, p provider.Provider) (metadata Metadata, diags diag.Diagnostics) {
	var providerMetadataResp provider.MetadataResponse
	p.Metadata(ctx, provider.MetadataRequest{}, &providerMetadataResp)

	metadata = Metadata{
		ProviderName: providerMetadataResp.TypeName,
		Resources:    ResourceMetadatas{},
		DataSources:  DataSourceMetadatas{},
	}

	for _, builder := range p.Resources(ctx) {
		res := builder()

		// Get the resource type
		var metadataResp resource.MetadataResponse
		res.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: providerMetadataResp.TypeName}, &metadataResp)
		resourceType := metadataResp.TypeName

		var schemaResp resource.SchemaResponse
		res.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
		diags.Append(schemaResp.Diagnostics...)
		if diags.HasError() {
			return
		}

		sch, odiags := NewResourceSchema(ctx, schemaResp.Schema)
		diags.Append(odiags...)
		if diags.HasError() {
			return
		}

		resMetadata := ResourceMetadata{
			Schema: sch,
		}

		if resourceWithIdentity, ok := res.(resource.ResourceWithIdentity); ok {
			var schemaResp resource.IdentitySchemaResponse
			resourceWithIdentity.IdentitySchema(ctx, resource.IdentitySchemaRequest{}, &schemaResp)
			diags.Append(schemaResp.Diagnostics...)
			if diags.HasError() {
				return
			}

			sch, odiags := NewResourceIdentitySchema(ctx, schemaResp.IdentitySchema)
			diags.Append(odiags...)
			if diags.HasError() {
				return
			}

			resMetadata.Identity = &sch
		}

		metadata.Resources[resourceType] = resMetadata
	}

	for _, builder := range p.DataSources(ctx) {
		ds := builder()

		// Get the data source type
		var metadataResp datasource.MetadataResponse
		ds.Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: providerMetadataResp.TypeName}, &metadataResp)
		datasourceType := metadataResp.TypeName

		var schemaResp datasource.SchemaResponse
		ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)
		diags.Append(schemaResp.Diagnostics...)
		if diags.HasError() {
			return
		}

		sch, odiags := NewDataSourceSchema(ctx, schemaResp.Schema)
		diags.Append(odiags...)
		if diags.HasError() {
			return
		}

		dsMetadata := DataSourceMetadata{
			Schema: sch,
		}

		metadata.DataSources[datasourceType] = dsMetadata
	}

	return
}
