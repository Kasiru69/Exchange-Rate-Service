package domain

import (
	"context"
	"time"
)

type ExchangeRepository interface {
	GetLatestRate(ctx context.Context, from, to string) (*ExchangeRate, error)
	GetHistoricalRate(ctx context.Context, from, to, date string) (*ExchangeRate, error)
	GetAllLatestRates(ctx context.Context, baseCurrency string) (map[string]float64, error)
}

type CacheRepository interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string, dest interface{}) error
	Delete(key string) error
	Clear() error
}
