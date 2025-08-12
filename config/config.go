package config

import "time"

const (
	// Upstream API URLs
	UpstreamGenerateURL = "https://gen.akash.network/api/generate"
	UpstreamStatusURL   = "https://gen.akash.network/api/status"
	UpstreamImageBase   = "https://gen.akash.network"

	// Timeout configurations
	GenerateTimeout    = 30 * time.Second
	StatusTimeout      = 10 * time.Second
	ImageFetchTimeout  = 30 * time.Second
	PollingInterval    = 1 * time.Second
	MaxPollingDuration = 5 * time.Minute

	// Concurrency configuration
	MaxConcurrentRequests = 10

	// Server configuration
	ServerPort = ":6571"
)

// PreferredGPUs returns the list of preferred GPU types
func PreferredGPUs() []string {
	return []string{
		"RTX4090",
		"A10",
		"A100",
		"V100-32Gi",
		"H100",
	}
}