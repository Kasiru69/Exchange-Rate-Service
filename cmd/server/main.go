package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"exchange-rate-service/internal/api"
	"exchange-rate-service/internal/config"
	"exchange-rate-service/internal/repository"
	"exchange-rate-service/internal/service"
	"exchange-rate-service/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	logger := logger.NewLogger(cfg.LogLevel)
	defer logger.Sync()

	cacheRepo := repository.NewCacheRepository()
	apiRepo := repository.NewExchangeAPIRepository(cfg.APIKey, cfg.BaseURL)

	exchangeService := service.NewExchangeService(cacheRepo, apiRepo, logger)

	go exchangeService.StartRateUpdater(context.Background())

	router := api.NewRouter(exchangeService, logger)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		logger.Info("Server starting on port " + cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server: " + err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: " + err.Error())
	}

	logger.Info("Server exited")
}
