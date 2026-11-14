package common

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type ErrorResponse struct {
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// DecodeResponse closes an HTTP response and decodes its JSON body into target.
func DecodeResponse(t *testing.T, response *http.Response, target interface{}) {
	t.Helper()
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "failed to read response body")
	if len(body) == 0 {
		return
	}
	require.NoErrorf(t, json.Unmarshal(body, target), "failed to decode response body: %s", string(body))
}
