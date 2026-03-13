package testprovider

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	ctx := context.Background()
	serveOpts := providerserver.ServeOpts{
		Address: "registry.terraform.io/magodo/examplecloud",
	}
	err := providerserver.Serve(ctx, func() provider.Provider { return &ExampleCloudProvider{} }, serveOpts)
	if err != nil {
		log.Fatalf("Error serving provider: %s", err)
	}
}
