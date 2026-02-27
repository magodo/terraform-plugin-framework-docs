package tfproviderdocs

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/magodo/terraform-provider-docs/internal/provider"
)

func TestNewProviderData(t *testing.T) {
	pd, diags := NewProviderData(t.Context(), &provider.ExampleCloudProvider{})
	if diags.HasError() {
		t.Fatal(diags.Errors())
	}

	spew.Dump(pd.Resources)
}
