package domain

import (
	"time"
)

type ExchangeRate struct {
	FromCurrency string    `json:"from_currency"`
	ToCurrency   string    `json:"to_currency"`
	Rate         float64   `json:"rate"`
	Timestamp    time.Time `json:"timestamp"`
	Date         string    `json:"date"`
}

type ConversionRequest struct {
	From   string  `json:"from" binding:"required"`
	To     string  `json:"to" binding:"required"`
	Amount float64 `json:"amount"`
	Date   string  `json:"date,omitempty"`
}

type ConversionResponse struct {
	Amount       float64   `json:"amount"`
	FromCurrency string    `json:"from_currency"`
	ToCurrency   string    `json:"to_currency"`
	Rate         float64   `json:"rate"`
	Date         string    `json:"date"`
	Timestamp    time.Time `json:"timestamp"`
}

type LatestRatesResponse struct {
	BaseCurrency string             `json:"base_currency"`
	Rates        map[string]float64 `json:"rates"`
	Timestamp    time.Time          `json:"timestamp"`
	Date         string             `json:"date"`
}

type HistoricalRatesResponse struct {
	FromCurrency string                  `json:"from_currency"`
	ToCurrency   string                  `json:"to_currency"`
	Rates        map[string]ExchangeRate `json:"rates"`
	StartDate    string                  `json:"start_date"`
	EndDate      string                  `json:"end_date"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var SupportedCurrencies = map[string]bool{
	"USD": true,
	"INR": true,
	"EUR": true,
	"JPY": true,
	"GBP": true,
}
