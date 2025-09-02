package utils

import (
	"fmt"
	"time"
)

func ValidateDate(dateStr string) error {
	if dateStr == "" {
		return fmt.Errorf("date cannot be empty")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format, expected YYYY-MM-DD")
	}

	now := time.Now()
	maxPastDate := now.AddDate(0, 0, -90)

	if date.After(now) {
		return fmt.Errorf("date cannot be in the future")
	}

	if date.Before(maxPastDate) {
		return fmt.Errorf("date cannot be more than 90 days in the past")
	}

	return nil
}

func GetDateRange(startDate, endDate string) ([]string, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format")
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format")
	}

	if start.After(end) {
		return nil, fmt.Errorf("start date cannot be after end date")
	}

	var dates []string
	current := start
	for !current.After(end) {
		dates = append(dates, current.Format("2006-01-02"))
		current = current.AddDate(0, 0, 1)
	}

	return dates, nil
}
