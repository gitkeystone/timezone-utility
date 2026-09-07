package lookup

import (
	"time"

	"timezone-utility/internal/clock"
	"timezone-utility/internal/data"
)

// Current computes the current local date/time in a timezone.
func Current(tzName string, loc *data.Location, c clock.Clock) (data.QueryResult, error) {
	locTZ, err := time.LoadLocation(tzName)
	if err != nil {
		return data.QueryResult{}, err
	}
	now := c.Now().In(locTZ)
	return data.QueryResult{
		ResolvedLocation: loc,
		LocalDate:        now.Format("2006-01-02"),
		LocalTime:        now.Format("15:04:05"),
		Timezone:         tzName,
		UTCOffset:        now.Format("-07:00"),
		DST:              now.IsDST(),
	}, nil
}
