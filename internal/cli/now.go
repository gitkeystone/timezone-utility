package cli

import (
	"timezone-utility/internal/data"
	"timezone-utility/internal/lookup"
	"timezone-utility/internal/output"
)

func runNow(args []string, app App) int {
	jsonOut, queries := extractJSON(args)
	if len(queries) == 0 {
		now := app.Clock.Now()
		res := data.QueryResult{
			Query:     "local",
			LocalDate: now.Format("2006-01-02"),
			LocalTime: now.Format("15:04:05"),
			Timezone:  "Local",
			UTCOffset: now.Format("-07:00"),
			DST:       now.IsDST(),
		}
		return emit([]data.QueryResult{res}, jsonOut, app)
	}

	results := make([]data.QueryResult, 0, len(queries))
	exitCode := 0
	for _, q := range queries {
		tzName, loc, err := lookup.Resolve(app.Store, q)
		if err != nil {
			output.WriteError(app.Err, err)
			exitCode = 2
			continue
		}
		res, err := lookup.Current(tzName, loc, app.Clock)
		if err != nil {
			output.WriteError(app.Err, err)
			exitCode = 2
			continue
		}
		res.Query = q
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
