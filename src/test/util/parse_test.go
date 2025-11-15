package util_test

import (
	"testing"

	"github.com/hellomyzn/nf-analysis/internal/util"
)

func Test_ConvertDate(t *testing.T) {
	input := "11/14/25"
	want := "2025-11-14"

	got, err := util.ConvertDate(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

