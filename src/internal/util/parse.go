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
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", "", ""
	}

	// シリーズものは区切りが2つ以上存在する想定。
	// 末尾2つの区切りを Season / Episode の区切りとして扱い、
	// それより前を Title として扱う。
	colonCount := strings.Count(trimmed, ":")
	if colonCount < 2 {
		return trimmed, "", ""
	}

	last := strings.LastIndex(trimmed, ":")
	if last == -1 {
		return trimmed, "", ""
	}

	episode := strings.TrimSpace(trimmed[last+1:])
	beforeEpisode := strings.TrimSpace(trimmed[:last])

	secondLast := strings.LastIndex(beforeEpisode, ":")
	if secondLast == -1 {
		return trimmed, "", ""
	}

	season := strings.TrimSpace(beforeEpisode[secondLast+1:])
	title := strings.TrimSpace(beforeEpisode[:secondLast])

	return title, season, episode
}
