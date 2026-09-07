package lookup

import (
	"testing"
	"time"

	"timezone-utility/internal/clock"
)

func TestCurrentShanghai(t *testing.T) {
	c := clock.FixedClock{T: time.Date(2026, 9, 8, 1, 57, 28, 0, time.UTC)}
	res, err := Current("Asia/Shanghai", nil, c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.LocalDate != "2026-09-08" {
		t.Errorf("date = %q, want 2026-09-08", res.LocalDate)
	}
	if res.LocalTime != "09:57:28" {
		t.Errorf("time = %q, want 09:57:28", res.LocalTime)
	}
	if res.UTCOffset != "+08:00" {
		t.Errorf("offset = %q, want +08:00", res.UTCOffset)
	}
	if res.Timezone != "Asia/Shanghai" {
		t.Errorf("timezone = %q", res.Timezone)
	}
}

func TestCurrentNewYorkDST(t *testing.T) {
	// September 2026: New York is on EDT (UTC-4).
	c := clock.FixedClock{T: time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC)}
	res, err := Current("America/New_York", nil, c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.UTCOffset != "-04:00" {
		t.Errorf("offset = %q, want -04:00", res.UTCOffset)
	}
	if !res.DST {
		t.Errorf("expected DST=true in September for New York")
	}
}
