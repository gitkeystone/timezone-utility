package convert

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

func TestConvertInstant(t *testing.T) {
	// Beijing 15:00 (UTC+8) = UTC 07:00.
	tBeijing := time.Date(2026, 9, 8, 15, 0, 0, 0, mustLoc("Asia/Shanghai"))
	res, err := ConvertInstant(tBeijing, "America/New_York", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// UTC 07:00 -> New York EDT (UTC-4) = 03:00.
	if res.LocalTime != "03:00" {
		t.Errorf("local time = %q, want 03:00", res.LocalTime)
	}
	if res.LocalDate != "2026-09-08" {
		t.Errorf("date = %q, want 2026-09-08", res.LocalDate)
	}
	if !res.DST {
		t.Errorf("expected DST=true in September for New York")
	}
	if res.UTCOffset != "-04:00" {
		t.Errorf("offset = %q, want -04:00", res.UTCOffset)
	}
}

func TestConvertInstantDateLine(t *testing.T) {
	// Auckland (UTC+12) 2026-09-09 01:00 = UTC 2026-09-08 13:00.
	tAuckland := time.Date(2026, 9, 9, 1, 0, 0, 0, mustLoc("Pacific/Auckland"))
	res, err := ConvertInstant(tAuckland, "America/Los_Angeles", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// UTC 13:00 -> Los Angeles PDT (UTC-7) = 06:00 on 2026-09-08.
	if res.LocalDate != "2026-09-08" {
		t.Errorf("date = %q, want 2026-09-08 (crossed date line)", res.LocalDate)
	}
	if res.LocalTime != "06:00" {
		t.Errorf("local time = %q, want 06:00", res.LocalTime)
	}
}
