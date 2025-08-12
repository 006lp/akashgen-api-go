package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/006lp/akashgen-api-go/config"
	"github.com/006lp/akashgen-api-go/handlers"
	"github.com/006lp/akashgen-api-go/middleware"
)

var (
	logger            *zap.Logger
	concurrentLimiter chan struct{}
	wg                sync.WaitGroup
)

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	concurrentLimiter = make(chan struct{}, config.MaxConcurrentRequests)
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(middleware.GinLogger(logger))
	r.Use(gin.Recovery())

	r.POST("/api/generate", handlers.HandleGenerate(logger, concurrentLimiter, &wg))
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	server := &http.Server{
		Addr:    ":6571",
		Handler: r,
	}

	go func() {
		logger.Info("Server starting", zap.String("addr", ":6571"))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	quit := make(chan struct{})
	go func() {
		time.Sleep(1 * time.Hour) // In production, use signal handling
		close(quit)
	}()

	<-quit
	logger.Info("Shutting down server...")

	close(concurrentLimiter)
	wg.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}