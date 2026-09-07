package integration

import (
	"encoding/json"
	"testing"

	"timezone-utility/internal/cli"
)

func TestJSONOutputContract(t *testing.T) {
	app, out, _ := newApp()
	code := cli.Run([]string{"now", "北京", "--json"}, app)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	var payload struct {
		Results []struct {
			Query     string `json:"query"`
			LocalDate string `json:"local_date"`
			LocalTime string `json:"local_time"`
			Timezone  string `json:"timezone"`
			UTCOffset string `json:"utc_offset"`
			DST       bool   `json:"dst"`
			Resolved  *struct {
				Name     string `json:"name"`
				Country  string `json:"country"`
				Timezone string `json:"timezone"`
			} `json:"resolved"`
		} `json:"results"`
	}
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal JSON: %v", err)
	}
	if len(payload.Results) != 1 {
		t.Fatalf("results = %d, want 1", len(payload.Results))
	}
	r := payload.Results[0]
	if r.Query != "北京" {
		t.Errorf("query = %q", r.Query)
	}
	if r.Timezone != "Asia/Shanghai" {
		t.Errorf("timezone = %q", r.Timezone)
	}
	if r.UTCOffset != "+08:00" {
		t.Errorf("utc_offset = %q", r.UTCOffset)
	}
	if r.LocalDate == "" || r.LocalTime == "" {
		t.Errorf("missing date/time: %+v", r)
	}
	if r.Resolved == nil || r.Resolved.Name != "Beijing" || r.Resolved.Timezone != "Asia/Shanghai" {
		t.Errorf("resolved = %+v", r.Resolved)
	}
}
