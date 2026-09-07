// Package convert converts a time instant between timezones.
package convert

import (
	"time"

	"timezone-utility/internal/data"
)

// ConvertInstant renders an instant in the target timezone as a QueryResult.
func ConvertInstant(t time.Time, tzName string, loc *data.Location) (data.QueryResult, error) {
	locTZ, err := time.LoadLocation(tzName)
	if err != nil {
		return data.QueryResult{}, err
	}
	lt := t.In(locTZ)
	return data.QueryResult{
		ResolvedLocation: loc,
		LocalDate:        lt.Format("2006-01-02"),
		LocalTime:        lt.Format("15:04"),
		Timezone:         tzName,
		UTCOffset:        lt.Format("-07:00"),
		DST:              lt.IsDST(),
	}, nil
}
