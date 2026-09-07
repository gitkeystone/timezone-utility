package cli

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"timezone-utility/internal/clock"
	"timezone-utility/internal/data"
)

func testApp(t *testing.T) (App, *bytes.Buffer, *bytes.Buffer) {
	t.Helper()
	store, err := data.Load()
	if err != nil {
		t.Fatalf("load store: %v", err)
	}
	var out, errBuf bytes.Buffer
	app := App{
		Store: store,
		Clock: clock.FixedClock{T: time.Date(2026, 9, 8, 1, 57, 28, 0, time.UTC)},
		Out:   &out,
		Err:   &errBuf,
	}
	return app, &out, &errBuf
}

func TestRunVersion(t *testing.T) {
	app, out, _ := testApp(t)
	if code := Run([]string{"--version"}, app); code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if !strings.Contains(out.String(), version) {
		t.Errorf("version output = %q", out.String())
	}
}

func TestRunHelp(t *testing.T) {
	app, out, _ := testApp(t)
	if code := Run([]string{"--help"}, app); code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if !strings.Contains(out.String(), "用法") {
		t.Errorf("help output = %q", out.String())
	}
}

func TestRunUnknownCommand(t *testing.T) {
	app, _, errBuf := testApp(t)
	if code := Run([]string{"frobnicate"}, app); code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}
	if !strings.Contains(errBuf.String(), "未知命令") {
		t.Errorf("stderr = %q", errBuf.String())
	}
}

func TestRunNoArgsLocalTime(t *testing.T) {
	app, out, _ := testApp(t)
	if code := Run(nil, app); code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if !strings.Contains(out.String(), "Local") {
		t.Errorf("output = %q", out.String())
	}
}

func TestExtractJSON(t *testing.T) {
	jsonOut, rest := extractJSON([]string{"北京", "--json", "上海"})
	if !jsonOut {
		t.Error("expected jsonOut=true")
	}
	if len(rest) != 2 || rest[0] != "北京" || rest[1] != "上海" {
		t.Errorf("rest = %v", rest)
	}
}

func TestParseTimeInSkippedTime(t *testing.T) {
	// 2026-03-08 02:30 in America/New_York falls in the spring-forward gap.
	if _, err := parseTimeIn("2026-03-08 02:30", "America/New_York"); err == nil {
		t.Error("expected error for non-existent DST time")
	}
}
