package cli

import (
	"flag"
	"fmt"

	"timezone-utility/internal/convert"
	"timezone-utility/internal/data"
	"timezone-utility/internal/lookup"
	"timezone-utility/internal/output"
)

func runConvert(args []string, app App) int {
	jsonOut, args := extractJSON(args)
	fs := flag.NewFlagSet("convert", flag.ContinueOnError)
	from := fs.String("from", "", "来源地点（时区/地名/邮编）")
	var tos multiFlag
	fs.Var(&tos, "to", "目标地点（可重复）")
	fs.SetOutput(app.Err)
	if err := fs.Parse(args); err != nil {
		return 1
	}
	if fs.NArg() != 1 {
		fmt.Fprintln(app.Err, "error: convert 需要一个时间参数（YYYY-MM-DD HH:mm 或 RFC3339）")
		return 1
	}
	if *from == "" || len(tos) == 0 {
		fmt.Fprintln(app.Err, "error: convert 需要 --from 与至少一个 --to")
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

	results := make([]data.QueryResult, 0, len(tos))
	exitCode := 0
	for _, target := range tos {
		tzName, loc, err := lookup.Resolve(app.Store, target)
		if err != nil {
			output.WriteError(app.Err, err)
			exitCode = 2
			continue
		}
		res, err := convert.ConvertInstant(t, tzName, loc)
		if err != nil {
			output.WriteError(app.Err, err)
			exitCode = 2
			continue
		}
		res.Query = target
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
