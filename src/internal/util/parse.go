package util

import (
	"fmt"
	"strings"
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

// SplitTitle 分割処理
// 例: "The Walking Dead: Season 5: Four Walls and a Roof"
// → title="The Walking Dead", season="Season 5", episode="Four Walls and a Roof"
func SplitTitle(input string) (string, string, string) {
	parts := strings.Split(input, ":")

	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	if len(parts) == 3 {
		return parts[0], parts[1], parts[2]
	}

	// 想定外の形式でもクラッシュしないようにフォールバック
	switch len(parts) {
	case 1:
		return parts[0], "", ""
	case 2:
		return parts[0], parts[0], ""
	default:
		// parts[0] = title, 残り全部 season or episode として結合
		return parts[0], parts[1], strings.Join(parts[2:], ":")
	}

}
