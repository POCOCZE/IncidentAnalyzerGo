package core

import "time"

// Parse RFC3339 compatible time string to a 'time.Time' data type.
func StringToTime(input string) (*time.Time, error) {
	output, err := time.Parse(time.RFC3339, input)
	if err != nil {
		return nil, err
	}
	return &output, nil
}

// Converts time to UTC
func ConvertToUTC(t *time.Time) *time.Time {
	utcTime := t.UTC()
	return &utcTime
}