// Package meeting provides cross-timezone meeting scheduling helpers.
package meeting

import (
	"time"

	"timezone-utility/internal/convert"
	"timezone-utility/internal/data"
)

// LocalTime renders an instant in a participant's timezone.
func LocalTime(t time.Time, tzName string, loc *data.Location) (data.QueryResult, error) {
	return convert.ConvertInstant(t, tzName, loc)
}
