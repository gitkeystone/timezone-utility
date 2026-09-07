package cli

import (
	"fmt"
	"strings"
	"time"
)

// multiFlag collects repeatable flag values.
type multiFlag []string

func (m *multiFlag) String() string { return strings.Join(*m, ",") }

func (m *multiFlag) Set(v string) error {
	*m = append(*m, v)
	return nil
}

// extractJSON removes any "--json" token from args (anywhere) and returns the
// json flag value along with the remaining args.
func extractJSON(args []string) (bool, []string) {
	jsonOut := false
	rest := make([]string, 0, len(args))
	for _, a := range args {
		if a == "--json" {
			jsonOut = true
			continue
		}
		rest = append(rest, a)
	}
	return jsonOut, rest
}

// parseTimeIn parses a time string in the given timezone.
func parseTimeIn(s, tzName string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	loc, err := time.LoadLocation(tzName)
	if err != nil {
		return time.Time{}, err
	}
	t, err := time.ParseInLocation("2006-01-02 15:04", s, loc)
	if err != nil {
		return time.Time{}, err
	}
	// Detect a non-existent local time (spring-forward gap): the wall clock is
	// normalized forward, so the parsed value no longer matches the input.
	if t.Format("2006-01-02 15:04") != s {
		return time.Time{}, fmt.Errorf("时间 %s 在时区 %s 中不存在（夏令时切换）", s, tzName)
	}
	return t, nil
}
