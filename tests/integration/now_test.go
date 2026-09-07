package integration

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"timezone-utility/internal/cli"
	"timezone-utility/internal/clock"
	"timezone-utility/internal/data"
)

func newApp() (cli.App, *bytes.Buffer, *bytes.Buffer) {
	store, err := data.Load()
	if err != nil {
		panic(err)
	}
	var out, errBuf bytes.Buffer
	app := cli.App{
		Store: store,
		Clock: clock.FixedClock{T: time.Date(2026, 9, 8, 1, 57, 28, 0, time.UTC)},
		Out:   &out,
		Err:   &errBuf,
	}
	return app, &out, &errBuf
}

func TestNowTimezone(t *testing.T) {
	app, out, _ := newApp()
	code := cli.Run([]string{"now", "Asia/Shanghai"}, app)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if !strings.Contains(out.String(), "2026-09-08") || !strings.Contains(out.String(), "09:57:28") {
		t.Errorf("output missing expected date/time: %q", out.String())
	}
	if !strings.Contains(out.String(), "+08:00") {
		t.Errorf("output missing offset: %q", out.String())
	}
}

func TestNowPlace(t *testing.T) {
	app, out, _ := newApp()
	code := cli.Run([]string{"now", "北京"}, app)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if !strings.Contains(out.String(), "Asia/Shanghai") {
		t.Errorf("output missing timezone: %q", out.String())
	}
}

func TestNowPostal(t *testing.T) {
	app, out, _ := newApp()
	code := cli.Run([]string{"now", "10001"}, app)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if !strings.Contains(out.String(), "America/New_York") {
		t.Errorf("output missing timezone: %q", out.String())
	}
}

func TestNowInvalid(t *testing.T) {
	app, _, errBuf := newApp()
	code := cli.Run([]string{"now", "Asia/Atlantis"}, app)
	if code != 2 {
		t.Fatalf("exit code = %d, want 2", code)
	}
	if !strings.Contains(errBuf.String(), "error:") {
		t.Errorf("stderr missing error: %q", errBuf.String())
	}
}

func TestNowMultiLocation(t *testing.T) {
	app, out, _ := newApp()
	code := cli.Run([]string{"now", "北京", "东京"}, app)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if !strings.Contains(out.String(), "Asia/Shanghai") || !strings.Contains(out.String(), "Asia/Tokyo") {
		t.Errorf("output missing timezones: %q", out.String())
	}
}
