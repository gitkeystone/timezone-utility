package meeting

import (
	"time"

	"timezone-utility/internal/data"
)

// ResolvedWindow is a participant availability window in a known timezone.
type ResolvedWindow struct {
	Location *data.Location
	Timezone string
	Loc      *time.Location
	Start    time.Time
	End      time.Time
}

// Overlap returns the common intersection of all windows, if any.
func Overlap(windows []ResolvedWindow) (start, end time.Time, ok bool) {
	if len(windows) == 0 {
		return time.Time{}, time.Time{}, false
	}
	start = windows[0].Start
	end = windows[0].End
	for _, w := range windows[1:] {
		if w.Start.After(start) {
			start = w.Start
		}
		if w.End.Before(end) {
			end = w.End
		}
	}
	if !start.Before(end) {
		return time.Time{}, time.Time{}, false
	}
	return start, end, true
}
