package lookup

import (
	"testing"
	"time"

	"timezone-utility/internal/clock"
	"timezone-utility/internal/data"
)

func BenchmarkCurrent(b *testing.B) {
	store, err := data.Load()
	if err != nil {
		b.Fatal(err)
	}
	c := clock.FixedClock{T: time.Date(2026, 9, 8, 1, 57, 28, 0, time.UTC)}
	tz, loc, err := Resolve(store, "北京")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Current(tz, loc, c); err != nil {
			b.Fatal(err)
		}
	}
}
