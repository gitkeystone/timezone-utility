package output

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"timezone-utility/internal/data"
)

func TestWriteText(t *testing.T) {
	var buf bytes.Buffer
	res := data.QueryResult{
		Query:     "北京",
		LocalDate: "2026-09-08",
		LocalTime: "10:00:00",
		Timezone:  "Asia/Shanghai",
		UTCOffset: "+08:00",
		DST:       false,
	}
	if err := WriteText(&buf, []data.QueryResult{res}); err != nil {
		t.Fatalf("WriteText: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"LOCATION", "2026-09-08", "10:00:00", "+08:00", "Asia/Shanghai"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q: %q", want, out)
		}
	}
}

func TestWriteJSON(t *testing.T) {
	var buf bytes.Buffer
	res := data.QueryResult{
		Query:     "北京",
		LocalDate: "2026-09-08",
		LocalTime: "10:00:00",
		Timezone:  "Asia/Shanghai",
		UTCOffset: "+08:00",
	}
	if err := WriteJSON(&buf, []data.QueryResult{res}); err != nil {
		t.Fatalf("WriteJSON: %v", err)
	}
	var out Results
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(out.Results) != 1 {
		t.Fatalf("results = %d, want 1", len(out.Results))
	}
	if out.Results[0].Timezone != "Asia/Shanghai" {
		t.Errorf("timezone = %q", out.Results[0].Timezone)
	}
}

func TestWriteError(t *testing.T) {
	var buf bytes.Buffer
	WriteError(&buf, errors.New("boom"))
	if !strings.HasPrefix(buf.String(), "error: ") {
		t.Errorf("got %q, want prefix \"error: \"", buf.String())
	}
}
