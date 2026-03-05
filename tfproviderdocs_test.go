package tfproviderdocs

import (
	"bytes"
	"os"
	"testing"

	"github.com/magodo/tfproviderdocs/internal/testprovider"
	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	g, err := NewGenerator(t.Context(), &testprovider.ExampleCloudProvider{})
	require.NoError(t, err)
	var buf bytes.Buffer
	require.NoError(t, g.RenderResource(t.Context(), &buf, "examplecloud_resource", nil))
	expected, err := os.ReadFile("./testdata/resource.md")
	require.NoError(t, err)
	require.Equal(t, string(expected), buf.String())
}
