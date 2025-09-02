package handlers

import (
	"net/http"
	"strconv"

	"exchange-rate-service/internal/domain"
	"exchange-rate-service/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ExchangeHandler struct {
	service *service.ExchangeService
	logger  *zap.Logger
}

func NewExchangeHandler(service *service.ExchangeService, logger *zap.Logger) *ExchangeHandler {
	return &ExchangeHandler{
		service: service,
		logger:  logger,
	}
}

func (h *ExchangeHandler) Convert(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	amountStr := c.Query("amount")
	date := c.Query("date")

	if from == "" || to == "" {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error:   "missing_parameters",
			Code:    400,
			Message: "from and to currencies are required",
		})
		return
	}

	amount := 1.0
	if amountStr != "" {
		var err error
		amount, err = strconv.ParseFloat(amountStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, domain.ErrorResponse{
				Error:   "invalid_amount",
				Code:    400,
				Message: "amount must be a valid number",
			})
			return
		}
	}

	req := &domain.ConversionRequest{
		From:   from,
		To:     to,
		Amount: amount,
		Date:   date,
	}

	result, err := h.service.ConvertCurrency(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("Conversion failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error:   "conversion_failed",
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *ExchangeHandler) GetLatestRates(c *gin.Context) {
	baseCurrency := c.Query("base")
	if baseCurrency == "" {
		baseCurrency = "USD"
	}

	result, err := h.service.GetLatestRates(c.Request.Context(), baseCurrency)
	if err != nil {
		h.logger.Error("Failed to get latest rates", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Error:   "fetch_failed",
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *ExchangeHandler) GetHistoricalRates(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if from == "" || to == "" || startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error:   "missing_parameters",
			Code:    400,
			Message: "from, to, start_date, and end_date are required",
		})
		return
	}

	result, err := h.service.GetHistoricalRates(c.Request.Context(), from, to, startDate, endDate)
	if err != nil {
		h.logger.Error("Failed to get historical rates", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Error:   "fetch_failed",
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *ExchangeHandler) GetSupportedCurrencies(c *gin.Context) {
	currencies := make([]string, 0, len(domain.SupportedCurrencies))
	for currency := range domain.SupportedCurrencies {
		currencies = append(currencies, currency)
	}

	c.JSON(http.StatusOK, gin.H{
		"currencies": currencies,
		"count":      len(currencies),
	})
}
