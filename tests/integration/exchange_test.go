package integration

import (
	//"bytes"
	//"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"exchange-rate-service/internal/api"
	//"exchange-rate-service/internal/domain"
	"exchange-rate-service/internal/repository"
	"exchange-rate-service/internal/service"

	//"exchange-rate-service/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestConvertEndpoint(t *testing.T) {
	logger := zap.NewNop()
	cacheRepo := repository.NewCacheRepository()
	apiRepo := repository.NewExchangeAPIRepository("", "")
	exchangeService := service.NewExchangeService(cacheRepo, apiRepo, logger)

	router := api.NewRouter(exchangeService, logger)

	tests := []struct {
		name           string
		url            string
		expectedStatus int
	}{
		{
			name:           "Valid conversion",
			url:            "/convert?from=USD&to=INR&amount=100",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Missing parameters",
			url:            "/convert?from=USD",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid currency",
			url:            "/convert?from=USD&to=XYZ&amount=100",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", tt.url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestHealthEndpoint(t *testing.T) {
	logger := zap.NewNop()
	cacheRepo := repository.NewCacheRepository()
	apiRepo := repository.NewExchangeAPIRepository("", "")
	exchangeService := service.NewExchangeService(cacheRepo, apiRepo, logger)

	router := api.NewRouter(exchangeService, logger)

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
}
