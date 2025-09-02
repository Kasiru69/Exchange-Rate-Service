package api

import (
	"exchange-rate-service/internal/api/handlers"
	"exchange-rate-service/internal/api/middleware"
	"exchange-rate-service/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewRouter(exchangeService *service.ExchangeService, logger *zap.Logger) *gin.Engine {
	router := gin.New()

	router.Use(middleware.Logger(logger))
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())

	exchangeHandler := handlers.NewExchangeHandler(exchangeService, logger)
	healthHandler := handlers.NewHealthHandler()

	router.GET("/health", healthHandler.Health)

	v1 := router.Group("/api/v1")
	{
		v1.GET("/convert", exchangeHandler.Convert)
		v1.GET("/latest", exchangeHandler.GetLatestRates)
		v1.GET("/historical", exchangeHandler.GetHistoricalRates)
		v1.GET("/currencies", exchangeHandler.GetSupportedCurrencies)
	}

	router.GET("/convert", exchangeHandler.Convert)

	return router
}
