package util

import (
	"fmt"
	"time"
)

func ConvertDate(input string) (string, error) {
	// Netflix の日付は "11/14/25" 形式（M/D/YY）
	t, err := time.Parse("1/2/06", input)
	if err != nil {
		return "", fmt.Errorf("failed to parse date: %w", err)
	}

	return t.Format("2006-01-02"), nil
}
