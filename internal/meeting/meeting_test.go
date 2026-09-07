package meeting

import (
	"testing"
	"time"
)

func mustLoc(name string) *time.Location {
	l, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return l
}

func TestLocalTime(t *testing.T) {
	// Beijing 2026-09-10 09:00 = UTC 01:00.
	tBeijing := time.Date(2026, 9, 10, 9, 0, 0, 0, mustLoc("Asia/Shanghai"))
	res, err := LocalTime(tBeijing, "Europe/London", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// UTC 01:00 -> London BST (UTC+1) = 02:00.
	if res.LocalTime != "02:00" {
		t.Errorf("local time = %q, want 02:00", res.LocalTime)
	}
	if res.LocalDate != "2026-09-10" {
		t.Errorf("date = %q, want 2026-09-10", res.LocalDate)
	}
}
