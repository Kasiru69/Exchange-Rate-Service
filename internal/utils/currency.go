package utils

import "exchange-rate-service/internal/domain"

func IsValidCurrency(currency string) bool {
	return domain.SupportedCurrencies[currency]
}
