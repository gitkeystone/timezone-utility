package integration

import (
	"strings"
	"testing"

	"timezone-utility/internal/cli"
)

func TestConvert(t *testing.T) {
	app, out, _ := newApp()
	code := cli.Run([]string{"convert", "--from", "北京", "--to", "纽约", "--to", "伦敦", "2026-09-08 15:00"}, app)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	// Beijing 15:00 = UTC 07:00 -> New York EDT 03:00, London BST 08:00.
	if !strings.Contains(out.String(), "03:00") || !strings.Contains(out.String(), "08:00") {
		t.Errorf("output missing expected times: %q", out.String())
	}
}
