package models

// GenerateRequest represents the incoming request for image generation
type GenerateRequest struct {
	Prompt    string `json:"prompt" binding:"required"`
	Negative  string `json:"negative"`
	Sampler   string `json:"sampler" binding:"required"`
	Scheduler string `json:"scheduler" binding:"required"`
}

// UpstreamGenerateRequest represents the request sent to upstream API
type UpstreamGenerateRequest struct {
	Prompt       string   `json:"prompt"`
	Negative     string   `json:"negative"`
	Sampler      string   `json:"sampler"`
	Scheduler    string   `json:"scheduler"`
	PreferredGpu []string `json:"preferred_gpu"`
}

// UpstreamStatusResponse represents the status response from upstream API
type UpstreamStatusResponse struct {
	JobID         string  `json:"job_id"`
	WorkerName    string  `json:"worker_name"`
	WorkerCity    string  `json:"worker_city"`
	WorkerCountry string  `json:"worker_country"`
	Status        string  `json:"status"`
	Result        string  `json:"result"`
	WorkerGPU     string  `json:"worker_gpu"`
	ElapsedTime   float64 `json:"elapsed_time"`
	QueuePosition int     `json:"queue_position"`
}

// GenerateJobResponse represents the response from generate request
type GenerateJobResponse struct {
	JobID string `json:"job_id"`
}