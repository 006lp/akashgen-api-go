package handlers

import (
	"context"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/006lp/akashgen-api-go/config"
	"github.com/006lp/akashgen-api-go/models"
	"github.com/006lp/akashgen-api-go/services"
)

// HandleGenerate handles the image generation request
func HandleGenerate(logger *zap.Logger, limiter chan struct{}, wg *sync.WaitGroup) gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter <- struct{}{}
		wg.Add(1)
		defer func() {
			<-limiter
			wg.Done()
		}()

		var req models.GenerateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Warn("Invalid JSON request", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid JSON",
				"details": err.Error(),
			})
			return
		}

		logger.Info("Received generate request",
			zap.String("prompt", req.Prompt),
			zap.String("sampler", req.Sampler),
			zap.String("scheduler", req.Scheduler),
		)

		upstreamReq := models.UpstreamGenerateRequest{
			Prompt:       req.Prompt,
			Negative:     req.Negative,
			Sampler:      req.Sampler,
			Scheduler:    req.Scheduler,
			PreferredGpu: config.PreferredGPUs(),
		}

		ctx, cancel := context.WithTimeout(context.Background(), config.GenerateTimeout)
		defer cancel()

		akashService := services.NewAkashService(logger)

		jobID, err := akashService.SendGenerateRequest(ctx, upstreamReq)
		if err != nil {
			logger.Error("Failed to generate image", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to generate image",
				"details": err.Error(),
			})
			return
		}

		logger.Info("Generate request sent successfully", zap.String("job_id", jobID))

		resultURL, err := akashService.PollJobStatus(ctx, jobID)
		if err != nil {
			logger.Error("Failed to poll job status", zap.Error(err), zap.String("job_id", jobID))
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to poll job status",
				"details": err.Error(),
			})
			return
		}

		imageService := services.NewImageService(logger)
		imageCtx, imageCancel := context.WithTimeout(context.Background(), config.ImageFetchTimeout)
		defer imageCancel()

		imageData, contentType, err := imageService.FetchImage(imageCtx, resultURL)
		if err != nil {
			logger.Error("Failed to fetch image", zap.Error(err), zap.String("result_url", resultURL))
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to fetch image",
				"details": err.Error(),
			})
			return
		}

		logger.Info("Image fetched successfully",
			zap.String("job_id", jobID),
			zap.String("content_type", contentType),
			zap.Int("size", len(imageData)),
		)

		c.Data(http.StatusOK, contentType, imageData)
	}
}
