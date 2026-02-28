package tfproviderdocs

import (
	"testing"

	"github.com/magodo/terraform-provider-docs/internal/provider"
)

func TestNewProviderData(t *testing.T) {
	pd, diags := NewProviderData(t.Context(), &provider.ExampleCloudProvider{})
	if diags.HasError() {
		t.Fatal(diags.Errors())
	}
	//spew.Dump(pd.Resources)

	render := ResourceRender{
		SchemaInfo: pd.Resources["examplecloud_resource"],
	}

	if err := render.Execute(t.Output()); err != nil {
		t.Fatal(err)
	}
}
