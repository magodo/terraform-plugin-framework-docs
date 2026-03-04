package schema

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type Metadata struct {
	ProviderName string
	Resources    ResourceSchemas
}

func GetMetadata(ctx context.Context, p provider.Provider) (metadata Metadata, diags diag.Diagnostics) {
	var providerMetadataResp provider.MetadataResponse
	p.Metadata(ctx, provider.MetadataRequest{}, &providerMetadataResp)

	metadata = Metadata{
		ProviderName: metadata.ProviderName,
		Resources:    ResourceSchemas{},
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

		info, odiags := NewResourceSchema(ctx, schemaResp.Schema)
		diags.Append(odiags...)
		if diags.HasError() {
			return
		}

		metadata.Resources[resourceType] = info
	}

	return
}
