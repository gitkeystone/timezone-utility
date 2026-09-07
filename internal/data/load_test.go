package data

import "testing"

func TestLoad(t *testing.T) {
	s, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(s.Locations) == 0 {
		t.Error("no locations loaded")
	}
	for _, loc := range s.Locations {
		if loc.Name == "Beijing" {
			found := false
			for _, p := range loc.PostalCodes {
				if p == "100000" {
					found = true
				}
			}
			if !found {
				t.Error("Beijing missing postal code 100000")
			}
		}
	}
}

func TestResolvePlace(t *testing.T) {
	s, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	loc, err := s.ResolvePlace("北京")
	if err != nil {
		t.Fatalf("ResolvePlace: %v", err)
	}
	if loc.Name != "Beijing" || loc.Timezone != "Asia/Shanghai" {
		t.Errorf("got %+v", loc)
	}
}

func TestResolvePlaceCaseInsensitive(t *testing.T) {
	s, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	loc, err := s.ResolvePlace("TOKYO")
	if err != nil {
		t.Fatalf("ResolvePlace: %v", err)
	}
	if loc.Name != "Tokyo" {
		t.Errorf("got %q", loc.Name)
	}
}

func TestResolvePlaceNotFound(t *testing.T) {
	s, _ := Load()
	if _, err := s.ResolvePlace("不存在的城市"); err == nil {
		t.Error("expected error for unknown place")
	}
}

func TestResolvePostal(t *testing.T) {
	s, _ := Load()
	loc, err := s.ResolvePostal("100000")
	if err != nil {
		t.Fatalf("ResolvePostal: %v", err)
	}
	if loc.Name != "Beijing" {
		t.Errorf("got %q", loc.Name)
	}
}

func TestResolvePostalNotFound(t *testing.T) {
	s, _ := Load()
	if _, err := s.ResolvePostal("999999"); err == nil {
		t.Error("expected error for unknown postal")
	}
}
