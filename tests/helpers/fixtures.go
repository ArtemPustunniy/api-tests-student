package helpers

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const onePixelPNG = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII="

// CreateTestPNG writes a valid one-pixel PNG into the test's temporary directory.
func CreateTestPNG(t *testing.T) string {
	t.Helper()

	content, err := base64.StdEncoding.DecodeString(onePixelPNG)
	require.NoError(t, err, "failed to decode the test image")

	filePath := filepath.Join(t.TempDir(), "photo.png")
	require.NoError(t, os.WriteFile(filePath, content, 0o600), "failed to create the test image")
	return filePath
}
