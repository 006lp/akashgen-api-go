package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/006lp/akashgen-api-go/config"
	"github.com/006lp/akashgen-api-go/models"
	"github.com/006lp/akashgen-api-go/utils"
)

// AkashService handles communication with Akash network API
type AkashService struct {
	logger *zap.Logger
	client *http.Client
}

// NewAkashService creates a new AkashService instance
func NewAkashService(logger *zap.Logger) *AkashService {
	return &AkashService{
		logger: logger,
		client: &http.Client{Timeout: config.GenerateTimeout},
	}
}

// SendGenerateRequest sends a generate request to upstream API
func (s *AkashService) SendGenerateRequest(ctx context.Context, req models.UpstreamGenerateRequest) (string, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", config.UpstreamGenerateURL, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upstream generate API error: %d, body: %s", resp.StatusCode, string(body))
	}

	var result models.GenerateJobResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.JobID, nil
}

// PollJobStatus polls the job status until completion
func (s *AkashService) PollJobStatus(ctx context.Context, jobID string) (string, error) {
	pollCtx, cancel := context.WithTimeout(ctx, config.MaxPollingDuration)
	defer cancel()

	ticker := time.NewTicker(config.PollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-pollCtx.Done():
			return "", fmt.Errorf("polling timeout exceeded: %w", pollCtx.Err())
		case <-ticker.C:
			status, err := s.checkJobStatus(ctx, jobID)
			if err != nil {
				s.logger.Warn("Failed to check job status", zap.Error(err), zap.String("job_id", jobID))
				continue
			}

			s.logger.Debug("Job status checked",
				zap.String("job_id", jobID),
				zap.String("status", status.Status),
				zap.Int("queue_position", status.QueuePosition),
			)

			switch status.Status {
			case "succeeded":
				return status.Result, nil
			case "failed", "cancelled", "timeout":
				return "", fmt.Errorf("job failed with status: %s", status.Status)
			case "pending", "waiting":
				// Continue polling
			default:
				s.logger.Warn("Unknown job status", zap.String("status", status.Status))
			}
		}
	}
}

func (s *AkashService) checkJobStatus(ctx context.Context, jobID string) (*models.UpstreamStatusResponse, error) {
	url := fmt.Sprintf("%s?ids=%s", config.UpstreamStatusURL, jobID)

	req, err := utils.NewHTTPRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create status request: %w", err)
	}

	client := &http.Client{Timeout: config.StatusTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to check status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status API error: %d, body: %s", resp.StatusCode, string(body))
	}

	var results []models.UpstreamStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode status response: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no status returned for job_id: %s", jobID)
	}

	return &results[0], nil
}