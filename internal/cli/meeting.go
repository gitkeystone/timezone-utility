package cli

import (
	"flag"
	"fmt"

	"timezone-utility/internal/data"
	"timezone-utility/internal/lookup"
	"timezone-utility/internal/meeting"
	"timezone-utility/internal/output"
)

func runMeeting(args []string, app App) int {
	jsonOut, args := extractJSON(args)
	fs := flag.NewFlagSet("meeting", flag.ContinueOnError)
	from := fs.String("from", "", "提议时间所在地点（时区/地名/邮编）")
	fs.SetOutput(app.Err)
	if err := fs.Parse(args); err != nil {
		return 1
	}
	if fs.NArg() < 2 {
		fmt.Fprintln(app.Err, "error: meeting 需要提议时间与至少一个参与者地点")
		return 1
	}
	if *from == "" {
		fmt.Fprintln(app.Err, "error: meeting 需要 --from 指定提议时间所在地点")
		return 1
	}

	fromTZ, _, err := lookup.Resolve(app.Store, *from)
	if err != nil {
		output.WriteError(app.Err, err)
		return 2
	}
	t, err := parseTimeIn(fs.Arg(0), fromTZ)
	if err != nil {
		output.WriteError(app.Err, fmt.Errorf("无效时间 %q: %w", fs.Arg(0), err))
		return 2
	}

	participants := fs.Args()[1:]
	results := make([]data.QueryResult, 0, len(participants))
	exitCode := 0
	for _, p := range participants {
		tzName, loc, err := lookup.Resolve(app.Store, p)
		if err != nil {
			output.WriteError(app.Err, err)
			exitCode = 2
			continue
		}
		res, err := meeting.LocalTime(t, tzName, loc)
		if err != nil {
			output.WriteError(app.Err, err)
			exitCode = 2
			continue
		}
		res.Query = p
		results = append(results, res)
	}
	if len(results) == 0 {
		return exitCode
	}
	if code := emit(results, jsonOut, app); code != 0 {
		return code
	}
	return exitCode
}
