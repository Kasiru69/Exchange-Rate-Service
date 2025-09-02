package unit

import (
	"context"
	"testing"
	"time"

	"exchange-rate-service/internal/domain"
	"exchange-rate-service/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockCacheRepository struct {
	mock.Mock
}

func (m *MockCacheRepository) Set(key string, value interface{}, expiration time.Duration) error {
	args := m.Called(key, value, expiration)
	return args.Error(0)
}

func (m *MockCacheRepository) Get(key string, dest interface{}) error {
	args := m.Called(key, dest)
	return args.Error(0)
}

func (m *MockCacheRepository) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockCacheRepository) Clear() error {
	args := m.Called()
	return args.Error(0)
}

type MockExchangeRepository struct {
	mock.Mock
}

func (m *MockExchangeRepository) GetLatestRate(ctx context.Context, from, to string) (*domain.ExchangeRate, error) {
	args := m.Called(ctx, from, to)
	return args.Get(0).(*domain.ExchangeRate), args.Error(1)
}

func (m *MockExchangeRepository) GetHistoricalRate(ctx context.Context, from, to, date string) (*domain.ExchangeRate, error) {
	args := m.Called(ctx, from, to, date)
	return args.Get(0).(*domain.ExchangeRate), args.Error(1)
}

func (m *MockExchangeRepository) GetAllLatestRates(ctx context.Context, baseCurrency string) (map[string]float64, error) {
	args := m.Called(ctx, baseCurrency)
	return args.Get(0).(map[string]float64), args.Error(1)
}

func TestExchangeService_ConvertCurrency(t *testing.T) {
	mockCache := new(MockCacheRepository)
	mockAPI := new(MockExchangeRepository)
	logger := zap.NewNop()

	service := service.NewExchangeService(mockCache, mockAPI, logger)

	mockCache.On("Get", mock.AnythingOfType("string"), mock.Anything).Return(assert.AnError)

	expectedRate := &domain.ExchangeRate{
		FromCurrency: "USD",
		ToCurrency:   "INR",
		Rate:         83.25,
		Timestamp:    time.Now(),
		Date:         "2025-09-01",
	}
	mockAPI.On("GetLatestRate", mock.Anything, "USD", "INR").Return(expectedRate, nil)

	mockCache.On("Set", mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("time.Duration")).Return(nil)

	req := &domain.ConversionRequest{
		From:   "USD",
		To:     "INR",
		Amount: 100,
	}

	result, err := service.ConvertCurrency(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 8325.0, result.Amount)
	assert.Equal(t, "USD", result.FromCurrency)
	assert.Equal(t, "INR", result.ToCurrency)
	assert.Equal(t, 83.25, result.Rate)

	mockCache.AssertExpectations(t)
	mockAPI.AssertExpectations(t)
}
