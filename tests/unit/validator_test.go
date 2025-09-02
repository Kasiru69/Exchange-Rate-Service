package unit

import (
	"testing"
	"time"

	"exchange-rate-service/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestValidateDate(t *testing.T) {
	tests := []struct {
		name    string
		date    string
		wantErr bool
	}{
		{"Valid date", "2025-08-01", false},
		{"Empty date", "", true},
		{"Invalid format", "01-08-2025", true},
		{"Future date", time.Now().AddDate(0, 0, 1).Format("2006-01-02"), true},
		{"Too old date", time.Now().AddDate(0, 0, -91).Format("2006-01-02"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.ValidateDate(tt.date)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsValidCurrency(t *testing.T) {
	tests := []struct {
		currency string
		valid    bool
	}{
		{"USD", true},
		{"INR", true},
		{"EUR", true},
		{"JPY", true},
		{"GBP", true},
		{"CAD", false},
		{"AUD", false},
	}

	for _, tt := range tests {
		t.Run(tt.currency, func(t *testing.T) {
			result := utils.IsValidCurrency(tt.currency)
			assert.Equal(t, tt.valid, result)
		})
	}
}
