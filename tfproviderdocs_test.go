package tfproviderdocs

import (
	"testing"

	"github.com/magodo/tfproviderdocs/internal/testprovider"
)

func TestRender(t *testing.T) {
	g, err := NewGenerator(t.Context(), &testprovider.ExampleCloudProvider{})
	if err != nil {
		t.Fatal(err)
	}
	if err := g.RenderResource(t.Context(), t.Output(), "examplecloud_resource", nil); err != nil {
		t.Fatal(err)
	}
}
