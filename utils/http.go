package utils

import (
	"context"
	"io"
	"net/http"
)

// NewHTTPRequest creates a new HTTP request with context
func NewHTTPRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, method, url, body)
}