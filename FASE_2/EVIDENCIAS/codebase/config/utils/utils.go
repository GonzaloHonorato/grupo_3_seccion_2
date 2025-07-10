package utils

import (
	"strings"
	"time"
)

func ParseDateTime(dateTimeStr string) (time.Time, error) {
	const layout = "2006-01-02T15:04:05Z"

	parsedTime, err := time.Parse(layout, dateTimeStr)

	if err != nil {
		return time.Time{}, err
	}

	return parsedTime.UTC(), nil
}


func ToCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "")
}
