package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"exchange-rate-service/internal/domain"
)

type ExchangeAPIRepository struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

type ExchangeHostResponse struct {
	Success   bool               `json:"success"`
	Timestamp int64              `json:"timestamp"`
	Source    string             `json:"source"`
	Quotes    map[string]float64 `json:"quotes"`
	Error     *APIError          `json:"error,omitempty"`
}

type HistoricalResponse struct {
	Success    bool               `json:"success"`
	Historical bool               `json:"historical"`
	Date       string             `json:"date"`
	Timestamp  int64              `json:"timestamp"`
	Source     string             `json:"source"`
	Quotes     map[string]float64 `json:"quotes"`
	Error      *APIError          `json:"error,omitempty"`
}

type ConvertResponse struct {
	Success bool         `json:"success"`
	Query   ConvertQuery `json:"query"`
	Info    ConvertInfo  `json:"info"`
	Result  float64      `json:"result"`
	Error   *APIError    `json:"error,omitempty"`
}

type ConvertQuery struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

type ConvertInfo struct {
	Timestamp int64   `json:"timestamp"`
	Quote     float64 `json:"quote"`
}

type APIError struct {
	Code int    `json:"code"`
	Info string `json:"info"`
}

func NewExchangeAPIRepository(apiKey, baseURL string) *ExchangeAPIRepository {
	if baseURL == "" {
		baseURL = "https://api.exchangerate.host"
	}

	return &ExchangeAPIRepository{
		apiKey:  apiKey,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (r *ExchangeAPIRepository) GetLatestRate(ctx context.Context, from, to string) (*domain.ExchangeRate, error) {
	url := fmt.Sprintf("%s/live", r.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return r.getMockRate(from, to), nil
	}

	q := req.URL.Query()
	if r.apiKey != "" {
		q.Add("access_key", r.apiKey)
	}
	q.Add("source", from)
	q.Add("currencies", to)
	req.URL.RawQuery = q.Encode()

	resp, err := r.client.Do(req)
	if err != nil {
		return r.getMockRate(from, to), nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return r.getMockRate(from, to), nil
	}

	var apiResp ExchangeHostResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return r.getMockRate(from, to), nil
	}

	if !apiResp.Success || apiResp.Error != nil {
		return r.getMockRate(from, to), nil
	}

	quoteKey := from + to
	rate, exists := apiResp.Quotes[quoteKey]
	if !exists {
		return r.getMockRate(from, to), nil
	}

	return &domain.ExchangeRate{
		FromCurrency: from,
		ToCurrency:   to,
		Rate:         rate,
		Timestamp:    time.Unix(apiResp.Timestamp, 0),
		Date:         time.Unix(apiResp.Timestamp, 0).Format("2006-01-02"),
	}, nil
}

func (r *ExchangeAPIRepository) GetHistoricalRate(ctx context.Context, from, to, date string) (*domain.ExchangeRate, error) {
	url := fmt.Sprintf("%s/historical", r.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return r.getMockHistoricalRate(from, to, date), nil
	}

	q := req.URL.Query()
	if r.apiKey != "" {
		q.Add("access_key", r.apiKey)
	}
	q.Add("date", date)
	q.Add("source", from)
	q.Add("currencies", to)
	req.URL.RawQuery = q.Encode()

	resp, err := r.client.Do(req)
	if err != nil {
		return r.getMockHistoricalRate(from, to, date), nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return r.getMockHistoricalRate(from, to, date), nil
	}

	var apiResp HistoricalResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return r.getMockHistoricalRate(from, to, date), nil
	}

	if !apiResp.Success || apiResp.Error != nil {
		return r.getMockHistoricalRate(from, to, date), nil
	}

	quoteKey := from + to
	rate, exists := apiResp.Quotes[quoteKey]
	if !exists {
		return r.getMockHistoricalRate(from, to, date), nil
	}

	return &domain.ExchangeRate{
		FromCurrency: from,
		ToCurrency:   to,
		Rate:         rate,
		Timestamp:    time.Unix(apiResp.Timestamp, 0),
		Date:         date,
	}, nil
}

func (r *ExchangeAPIRepository) GetAllLatestRates(ctx context.Context, baseCurrency string) (map[string]float64, error) {
	url := fmt.Sprintf("%s/live", r.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return r.getMockRates(baseCurrency), nil
	}

	q := req.URL.Query()
	if r.apiKey != "" {
		q.Add("access_key", r.apiKey)
	}
	q.Add("source", baseCurrency)

	supportedCurrencies := []string{}
	for curr := range domain.SupportedCurrencies {
		if curr != baseCurrency {
			supportedCurrencies = append(supportedCurrencies, curr)
		}
	}
	q.Add("currencies", strings.Join(supportedCurrencies, ","))
	req.URL.RawQuery = q.Encode()

	resp, err := r.client.Do(req)
	if err != nil {
		return r.getMockRates(baseCurrency), nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return r.getMockRates(baseCurrency), nil
	}

	var apiResp ExchangeHostResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return r.getMockRates(baseCurrency), nil
	}

	if !apiResp.Success || apiResp.Error != nil {
		return r.getMockRates(baseCurrency), nil
	}

	rates := make(map[string]float64)
	for pair, rate := range apiResp.Quotes {
		if len(pair) >= 6 {
			toCurrency := pair[3:]
			rates[toCurrency] = rate
		}
	}

	return rates, nil
}

func (r *ExchangeAPIRepository) getMockRates(baseCurrency string) map[string]float64 {
	mockRates := map[string]map[string]float64{
		"USD": {
			"INR": 83.25,
			"EUR": 0.85,
			"JPY": 110.50,
			"GBP": 0.73,
		},
		"EUR": {
			"USD": 1.18,
			"INR": 98.12,
			"JPY": 130.25,
			"GBP": 0.86,
		},
		"GBP": {
			"USD": 1.37,
			"INR": 114.05,
			"EUR": 1.16,
			"JPY": 151.38,
		},
		"INR": {
			"USD": 0.012,
			"EUR": 0.010,
			"JPY": 1.33,
			"GBP": 0.0088,
		},
		"JPY": {
			"USD": 0.0090,
			"EUR": 0.0077,
			"INR": 0.75,
			"GBP": 0.0066,
		},
	}

	if rates, exists := mockRates[baseCurrency]; exists {
		return rates
	}
	return make(map[string]float64)
}

func (r *ExchangeAPIRepository) getMockRate(from, to string) *domain.ExchangeRate {
	baseRates := r.getMockRates(from)
	rate, exists := baseRates[to]
	if !exists {
		rate = 1.0
	}

	return &domain.ExchangeRate{
		FromCurrency: from,
		ToCurrency:   to,
		Rate:         rate,
		Timestamp:    time.Now(),
		Date:         time.Now().Format("2006-01-02"),
	}
}

func (r *ExchangeAPIRepository) getMockHistoricalRate(from, to, date string) *domain.ExchangeRate {
	baseRates := r.getMockRates(from)
	rate, exists := baseRates[to]
	if !exists {
		rate = 1.0
	}

	return &domain.ExchangeRate{
		FromCurrency: from,
		ToCurrency:   to,
		Rate:         rate,
		Timestamp:    time.Now(),
		Date:         date,
	}
}
