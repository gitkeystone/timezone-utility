package meeting

import (
	"testing"
	"time"
)

func TestOverlap(t *testing.T) {
	bj := mustLoc("Asia/Shanghai")
	lon := mustLoc("Europe/London")
	windows := []ResolvedWindow{
		{
			Timezone: "Asia/Shanghai",
			Loc:      bj,
			Start:    time.Date(2026, 9, 10, 16, 0, 0, 0, bj), // UTC 08:00
			End:      time.Date(2026, 9, 10, 20, 0, 0, 0, bj), // UTC 12:00
		},
		{
			Timezone: "Europe/London",
			Loc:      lon,
			Start:    time.Date(2026, 9, 10, 9, 0, 0, 0, lon),  // UTC 08:00
			End:      time.Date(2026, 9, 10, 12, 0, 0, 0, lon), // UTC 11:00
		},
	}
	start, end, ok := Overlap(windows)
	if !ok {
		t.Fatal("expected overlap")
	}
	wantStart := time.Date(2026, 9, 10, 8, 0, 0, 0, time.UTC)
	wantEnd := time.Date(2026, 9, 10, 11, 0, 0, 0, time.UTC)
	if !start.Equal(wantStart) {
		t.Errorf("start = %v, want %v", start.UTC(), wantStart)
	}
	if !end.Equal(wantEnd) {
		t.Errorf("end = %v, want %v", end.UTC(), wantEnd)
	}
}

func TestOverlapNone(t *testing.T) {
	bj := mustLoc("Asia/Shanghai")
	lon := mustLoc("Europe/London")
	windows := []ResolvedWindow{
		{
			Timezone: "Asia/Shanghai",
			Loc:      bj,
			Start:    time.Date(2026, 9, 10, 9, 0, 0, 0, bj),  // UTC 01:00
			End:      time.Date(2026, 9, 10, 12, 0, 0, 0, bj), // UTC 04:00
		},
		{
			Timezone: "Europe/London",
			Loc:      lon,
			Start:    time.Date(2026, 9, 10, 9, 0, 0, 0, lon),  // UTC 08:00
			End:      time.Date(2026, 9, 10, 12, 0, 0, 0, lon), // UTC 11:00
		},
	}
	if _, _, ok := Overlap(windows); ok {
		t.Error("expected no overlap")
	}
}
