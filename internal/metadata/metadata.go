package metadata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type Metadata struct {
	ProviderName   string
	ProviderSchema ProviderSchema
	Resources      ResourceMetadatas
	DataSources    DataSourceMetadatas
	Ephemerals     EphemeralMetadatas
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

type EphemeralMetadatas map[string]EphemeralMetadata

type EphemeralMetadata struct {
	Schema EphemeralSchema
}

func GetMetadata(ctx context.Context, p provider.Provider) (metadata Metadata, diags diag.Diagnostics) {
	var providerMetadataResp provider.MetadataResponse
	p.Metadata(ctx, provider.MetadataRequest{}, &providerMetadataResp)

	var schemaResp provider.SchemaResponse
	p.Schema(ctx, provider.SchemaRequest{}, &schemaResp)
	diags.Append(schemaResp.Diagnostics...)
	if diags.HasError() {
		return
	}
	providerSchema, odiags := NewProviderSchema(ctx, schemaResp.Schema)
	diags.Append(odiags...)
	if diags.HasError() {
		return
	}

	metadata = Metadata{
		ProviderName:   providerMetadataResp.TypeName,
		ProviderSchema: providerSchema,
		Resources:      ResourceMetadatas{},
		DataSources:    DataSourceMetadatas{},
		Ephemerals:     EphemeralMetadatas{},
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

		// Get the resource type
		var metadataResp datasource.MetadataResponse
		ds.Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: providerMetadataResp.TypeName}, &metadataResp)
		resourceType := metadataResp.TypeName

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

		metadata.DataSources[resourceType] = dsMetadata
	}

	if p, ok := p.(provider.ProviderWithEphemeralResources); ok {
		for _, builder := range p.EphemeralResources(ctx) {
			er := builder()

			// Get the resource type
			var metadataResp ephemeral.MetadataResponse
			er.Metadata(ctx, ephemeral.MetadataRequest{ProviderTypeName: providerMetadataResp.TypeName}, &metadataResp)
			resourceType := metadataResp.TypeName

			var schemaResp ephemeral.SchemaResponse
			er.Schema(ctx, ephemeral.SchemaRequest{}, &schemaResp)
			diags.Append(schemaResp.Diagnostics...)
			if diags.HasError() {
				return
			}

			sch, odiags := NewEphemeralSchema(ctx, schemaResp.Schema)
			diags.Append(odiags...)
			if diags.HasError() {
				return
			}

			emetadata := EphemeralMetadata{
				Schema: sch,
			}

			metadata.Ephemerals[resourceType] = emetadata
		}
	}

	return
}
