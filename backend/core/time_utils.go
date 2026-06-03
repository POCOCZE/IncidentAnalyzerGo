package core

import "time"

func StringToTime(input string) (*time.Time, error) {
	output, err := time.Parse(time.RFC3339, input)
	if err != nil {
		return nil, err
	}
	return &output, nil
}

func ConvertToUTC(t *time.Time) *time.Time {
	// Converts time to UTC
	utcTime := t.UTC()
	return &utcTime
}