package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strings"
	"time"

	"timezone-utility/internal/data"
	"timezone-utility/internal/lookup"
	"timezone-utility/internal/meeting"
	"timezone-utility/internal/output"
)

func runOverlap(args []string, app App) int {
	jsonOut, args := extractJSON(args)
	fs := flag.NewFlagSet("overlap", flag.ContinueOnError)
	var windows multiFlag
	fs.Var(&windows, "window", "参与者的可用时段（loc=HH:mm-HH:mm，可重复）")
	fs.SetOutput(app.Err)
	if err := fs.Parse(args); err != nil {
		return 1
	}
	if len(windows) < 2 {
		fmt.Fprintln(app.Err, "error: overlap 需要至少两个 --window")
		return 1
	}

	reference := app.Clock.Now()
	resolved := make([]meeting.ResolvedWindow, 0, len(windows))
	for _, w := range windows {
		rw, err := buildResolvedWindow(w, reference, app.Store)
		if err != nil {
			output.WriteError(app.Err, err)
			return 2
		}
		resolved = append(resolved, rw)
	}

	start, end, ok := meeting.Overlap(resolved)
	if !ok {
		if jsonOut {
			fmt.Fprintln(app.Out, `{"overlap": null}`)
		} else {
			fmt.Fprintln(app.Out, "无重叠时段")
		}
		return 0
	}
	if jsonOut {
		if err := writeOverlapJSON(app.Out, resolved, start, end); err != nil {
			output.WriteError(app.Err, err)
			return 2
		}
		return 0
	}
	writeOverlapText(app.Out, resolved, start, end)
	return 0
}

func buildResolvedWindow(w string, reference time.Time, store *data.Store) (meeting.ResolvedWindow, error) {
	locName, startStr, endStr, err := parseWindowString(w)
	if err != nil {
		return meeting.ResolvedWindow{}, err
	}
	tzName, loc, err := lookup.Resolve(store, locName)
	if err != nil {
		return meeting.ResolvedWindow{}, err
	}
	locTZ, err := time.LoadLocation(tzName)
	if err != nil {
		return meeting.ResolvedWindow{}, err
	}
	startHM, err := time.ParseInLocation("15:04", startStr, locTZ)
	if err != nil {
		return meeting.ResolvedWindow{}, fmt.Errorf("无效开始时间 %q: %w", startStr, err)
	}
	endHM, err := time.ParseInLocation("15:04", endStr, locTZ)
	if err != nil {
		return meeting.ResolvedWindow{}, fmt.Errorf("无效结束时间 %q: %w", endStr, err)
	}
	start := time.Date(reference.Year(), reference.Month(), reference.Day(), startHM.Hour(), startHM.Minute(), 0, 0, locTZ)
	end := time.Date(reference.Year(), reference.Month(), reference.Day(), endHM.Hour(), endHM.Minute(), 0, 0, locTZ)
	if !start.Before(end) {
		return meeting.ResolvedWindow{}, fmt.Errorf("窗口开始时间必须早于结束时间: %q", w)
	}
	return meeting.ResolvedWindow{Location: loc, Timezone: tzName, Loc: locTZ, Start: start, End: end}, nil
}

func writeOverlapText(w io.Writer, resolved []meeting.ResolvedWindow, start, end time.Time) {
	fmt.Fprintln(w, "重叠时段:")
	fmt.Fprintf(w, "  %-12s %s - %s\n", "UTC", start.UTC().Format("15:04"), end.UTC().Format("15:04"))
	for _, rw := range resolved {
		name := rw.Timezone
		if rw.Location != nil {
			name = rw.Location.Name
		}
		fmt.Fprintf(w, "  %-12s %s - %s (%s)\n", name, start.In(rw.Loc).Format("15:04"), end.In(rw.Loc).Format("15:04"), rw.Timezone)
	}
}

func writeOverlapJSON(w io.Writer, resolved []meeting.ResolvedWindow, start, end time.Time) error {
	type participant struct {
		Location string `json:"location"`
		Timezone string `json:"timezone"`
		Start    string `json:"start"`
		End      string `json:"end"`
	}
	type overlapOut struct {
		StartUTC     string        `json:"start_utc"`
		EndUTC       string        `json:"end_utc"`
		Participants []participant `json:"participants"`
	}
	out := overlapOut{
		StartUTC: start.UTC().Format("15:04"),
		EndUTC:   end.UTC().Format("15:04"),
	}
	for _, rw := range resolved {
		name := rw.Timezone
		if rw.Location != nil {
			name = rw.Location.Name
		}
		out.Participants = append(out.Participants, participant{
			Location: name,
			Timezone: rw.Timezone,
			Start:    start.In(rw.Loc).Format("15:04"),
			End:      end.In(rw.Loc).Format("15:04"),
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func parseWindowString(w string) (string, string, string, error) {
	parts := strings.SplitN(w, "=", 2)
	if len(parts) != 2 {
		return "", "", "", fmt.Errorf("无效窗口 %q（应为 loc=HH:mm-HH:mm）", w)
	}
	times := strings.SplitN(parts[1], "-", 2)
	if len(times) != 2 {
		return "", "", "", fmt.Errorf("无效窗口 %q（应为 loc=HH:mm-HH:mm）", w)
	}
	return parts[0], times[0], times[1], nil
}
