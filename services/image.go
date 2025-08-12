package services

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/006lp/akashgen-api-go/config"
	"github.com/006lp/akashgen-api-go/utils"
)

// ImageService handles image fetching operations
type ImageService struct {
	logger *zap.Logger
	client *http.Client
}

// NewImageService creates a new ImageService instance
func NewImageService(logger *zap.Logger) *ImageService {
	return &ImageService{
		logger: logger,
		client: &http.Client{Timeout: config.ImageFetchTimeout},
	}
}

// FetchImage fetches the generated image from the upstream server
func (s *ImageService) FetchImage(ctx context.Context, imagePath string) ([]byte, string, error) {
	url := config.UpstreamImageBase + imagePath + "&w=2048&q=100"

	req, err := utils.NewHTTPRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create image request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("image fetch error: %d, body: %s", resp.StatusCode, string(body))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read image data: %w", err)
	}

	contentType := resp.Header.Get("Content-Type")
	return data, contentType, nil
}