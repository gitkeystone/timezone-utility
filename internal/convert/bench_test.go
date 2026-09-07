package convert

import (
	"testing"
	"time"
)

func BenchmarkConvertInstant(b *testing.B) {
	t := time.Date(2026, 9, 8, 15, 0, 0, 0, mustLoc("Asia/Shanghai"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := ConvertInstant(t, "America/New_York", nil); err != nil {
			b.Fatal(err)
		}
	}
}
