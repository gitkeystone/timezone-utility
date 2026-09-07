package integration

import (
	"strings"
	"testing"

	"timezone-utility/internal/cli"
)

func TestMeeting(t *testing.T) {
	app, out, _ := newApp()
	code := cli.Run([]string{"meeting", "--from", "北京", "2026-09-10 09:00", "伦敦", "纽约"}, app)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	// Beijing 09:00 = UTC 01:00 -> London BST 02:00 (09-10), New York EDT 21:00 (09-09).
	if !strings.Contains(out.String(), "02:00") || !strings.Contains(out.String(), "21:00") {
		t.Errorf("output missing expected times: %q", out.String())
	}
	if !strings.Contains(out.String(), "2026-09-09") {
		t.Errorf("output missing crossed-date for New York: %q", out.String())
	}
}
