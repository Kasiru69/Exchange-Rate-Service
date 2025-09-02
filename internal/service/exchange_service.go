package service

import (
	"context"
	"fmt"
	"time"

	"exchange-rate-service/internal/domain"
	"exchange-rate-service/internal/utils"

	"go.uber.org/zap"
)

type ExchangeService struct {
	cacheRepo domain.CacheRepository
	apiRepo   domain.ExchangeRepository
	logger    *zap.Logger
}

func NewExchangeService(cacheRepo domain.CacheRepository, apiRepo domain.ExchangeRepository, logger *zap.Logger) *ExchangeService {
	return &ExchangeService{
		cacheRepo: cacheRepo,
		apiRepo:   apiRepo,
		logger:    logger,
	}
}

func (s *ExchangeService) ConvertCurrency(ctx context.Context, req *domain.ConversionRequest) (*domain.ConversionResponse, error) {
	if !utils.IsValidCurrency(req.From) || !utils.IsValidCurrency(req.To) {
		return nil, fmt.Errorf("unsupported currency pair: %s to %s", req.From, req.To)
	}

	if req.Amount == 0 {
		req.Amount = 1
	}

	var rate *domain.ExchangeRate
	var err error

	if req.Date != "" {
		if err := utils.ValidateDate(req.Date); err != nil {
			return nil, err
		}

		rate, err = s.getHistoricalRate(ctx, req.From, req.To, req.Date)
	} else {
		rate, err = s.getLatestRate(ctx, req.From, req.To)
	}

	if err != nil {
		return nil, err
	}

	convertedAmount := req.Amount * rate.Rate

	return &domain.ConversionResponse{
		Amount:       convertedAmount,
		FromCurrency: req.From,
		ToCurrency:   req.To,
		Rate:         rate.Rate,
		Date:         rate.Date,
		Timestamp:    rate.Timestamp,
	}, nil
}

func (s *ExchangeService) GetLatestRates(ctx context.Context, baseCurrency string) (*domain.LatestRatesResponse, error) {
	if !utils.IsValidCurrency(baseCurrency) {
		return nil, fmt.Errorf("unsupported base currency: %s", baseCurrency)
	}

	cacheKey := fmt.Sprintf("latest_rates_%s", baseCurrency)

	var cachedRates map[string]float64
	if err := s.cacheRepo.Get(cacheKey, &cachedRates); err == nil {
		return &domain.LatestRatesResponse{
			BaseCurrency: baseCurrency,
			Rates:        cachedRates,
			Timestamp:    time.Now(),
			Date:         time.Now().Format("2006-01-02"),
		}, nil
	}

	rates, err := s.apiRepo.GetAllLatestRates(ctx, baseCurrency)
	if err != nil {
		return nil, err
	}

	s.cacheRepo.Set(cacheKey, rates, time.Hour)

	return &domain.LatestRatesResponse{
		BaseCurrency: baseCurrency,
		Rates:        rates,
		Timestamp:    time.Now(),
		Date:         time.Now().Format("2006-01-02"),
	}, nil
}

func (s *ExchangeService) GetHistoricalRates(ctx context.Context, from, to, startDate, endDate string) (*domain.HistoricalRatesResponse, error) {
	if !utils.IsValidCurrency(from) || !utils.IsValidCurrency(to) {
		return nil, fmt.Errorf("unsupported currency pair: %s to %s", from, to)
	}

	if err := utils.ValidateDate(startDate); err != nil {
		return nil, fmt.Errorf("invalid start date: %v", err)
	}

	if err := utils.ValidateDate(endDate); err != nil {
		return nil, fmt.Errorf("invalid end date: %v", err)
	}

	dates, err := utils.GetDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	rates := make(map[string]domain.ExchangeRate)
	for _, date := range dates {
		rate, err := s.getHistoricalRate(ctx, from, to, date)
		if err != nil {
			s.logger.Warn("Failed to get historical rate",
				zap.String("date", date),
				zap.Error(err))
			continue
		}
		rates[date] = *rate
	}

	return &domain.HistoricalRatesResponse{
		FromCurrency: from,
		ToCurrency:   to,
		Rates:        rates,
		StartDate:    startDate,
		EndDate:      endDate,
	}, nil
}

func (s *ExchangeService) getLatestRate(ctx context.Context, from, to string) (*domain.ExchangeRate, error) {
	cacheKey := fmt.Sprintf("rate_%s_%s_latest", from, to)

	var cachedRate domain.ExchangeRate
	if err := s.cacheRepo.Get(cacheKey, &cachedRate); err == nil {
		return &cachedRate, nil
	}

	rate, err := s.apiRepo.GetLatestRate(ctx, from, to)
	if err != nil {
		return nil, err
	}

	s.cacheRepo.Set(cacheKey, rate, time.Hour)

	return rate, nil
}

func (s *ExchangeService) getHistoricalRate(ctx context.Context, from, to, date string) (*domain.ExchangeRate, error) {
	cacheKey := fmt.Sprintf("rate_%s_%s_%s", from, to, date)

	var cachedRate domain.ExchangeRate
	if err := s.cacheRepo.Get(cacheKey, &cachedRate); err == nil {
		return &cachedRate, nil
	}

	rate, err := s.apiRepo.GetHistoricalRate(ctx, from, to, date)
	if err != nil {
		return nil, err
	}

	s.cacheRepo.Set(cacheKey, rate, 24*time.Hour)

	return rate, nil
}

func (s *ExchangeService) StartRateUpdater(ctx context.Context) {
	ticker := time.NewTicker(4 * time.Hour)
	defer ticker.Stop()

	s.updateAllRates(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.updateAllRates(ctx)
		}
	}
}

func (s *ExchangeService) updateAllRates(ctx context.Context) {
	currencies := []string{"USD", "EUR", "GBP", "INR", "JPY"}

	for _, baseCurrency := range currencies {
		rates, err := s.apiRepo.GetAllLatestRates(ctx, baseCurrency)
		if err != nil {
			s.logger.Error("Failed to update rates",
				zap.String("base_currency", baseCurrency),
				zap.Error(err))
			continue
		}

		cacheKey := fmt.Sprintf("latest_rates_%s", baseCurrency)
		s.cacheRepo.Set(cacheKey, rates, time.Hour)

		s.logger.Info("Updated rates", zap.String("base_currency", baseCurrency))
	}
}
