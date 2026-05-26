package main

import "time"

func IsValidTime(t *time.Time) bool {
	if t == nil || t.IsZero() {
		return false
	}
	return true
}

func TimeToString(t *time.Time) string {
	return t.Format(time.RFC3339)
}

func StringToTime(input string) (*time.Time, error) {
	output, err := time.Parse(time.RFC3339, input)
	if err != nil {
		return nil, err
	}
	return &output, nil
}

func ConvertToUTC(t *time.Time) *time.Time {
	// Converts time to UTC
	var utcTime time.Time
	utcTime = t.UTC()
	return &utcTime
}