package lookup

import (
	"testing"

	"timezone-utility/internal/data"
)

func testStore(t *testing.T) *data.Store {
	t.Helper()
	s, err := data.Load()
	if err != nil {
		t.Fatalf("load store: %v", err)
	}
	return s
}

func TestResolveTimezone(t *testing.T) {
	s := testStore(t)
	tz, loc, err := Resolve(s, "Asia/Shanghai")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tz != "Asia/Shanghai" || loc != nil {
		t.Errorf("got tz=%q loc=%v", tz, loc)
	}
}

func TestResolvePlaceAlias(t *testing.T) {
	s := testStore(t)
	tz, loc, err := Resolve(s, "北京")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tz != "Asia/Shanghai" || loc == nil || loc.Name != "Beijing" {
		t.Errorf("got tz=%q loc=%+v", tz, loc)
	}
}

func TestResolvePlaceEnglish(t *testing.T) {
	s := testStore(t)
	tz, loc, err := Resolve(s, "Tokyo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tz != "Asia/Tokyo" || loc == nil || loc.Name != "Tokyo" {
		t.Errorf("got tz=%q loc=%+v", tz, loc)
	}
}

func TestResolvePostal(t *testing.T) {
	s := testStore(t)
	tz, loc, err := Resolve(s, "100000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tz != "Asia/Shanghai" || loc == nil || loc.Name != "Beijing" {
		t.Errorf("got tz=%q loc=%+v", tz, loc)
	}
}

func TestResolvePostalUS(t *testing.T) {
	s := testStore(t)
	tz, loc, err := Resolve(s, "10001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tz != "America/New_York" || loc == nil || loc.Name != "New York" {
		t.Errorf("got tz=%q loc=%+v", tz, loc)
	}
}

func TestResolveInvalid(t *testing.T) {
	s := testStore(t)
	if _, _, err := Resolve(s, "Asia/Atlantis"); err == nil {
		t.Error("expected error for invalid timezone")
	}
	if _, _, err := Resolve(s, "不存在的城市"); err == nil {
		t.Error("expected error for unknown place")
	}
	if _, _, err := Resolve(s, "999999"); err == nil {
		t.Error("expected error for unknown postal")
	}
}
