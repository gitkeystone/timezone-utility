package integration

import (
	"strings"
	"testing"

	"timezone-utility/internal/cli"
)

func TestOverlap(t *testing.T) {
	app, out, _ := newApp()
	code := cli.Run([]string{"overlap", "--window", "北京=16:00-20:00", "--window", "伦敦=09:00-12:00"}, app)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	// Beijing 16:00-20:00 = UTC 08:00-12:00; London 09:00-12:00 BST = UTC 08:00-11:00.
	// Overlap = UTC 08:00-11:00 -> Beijing 16:00-19:00, London 09:00-12:00.
	if !strings.Contains(out.String(), "重叠时段") {
		t.Errorf("output missing overlap header: %q", out.String())
	}
	if !strings.Contains(out.String(), "16:00") || !strings.Contains(out.String(), "09:00") {
		t.Errorf("output missing local times: %q", out.String())
	}
}
