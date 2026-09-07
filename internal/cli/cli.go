// Package cli implements the command-line interface dispatch and subcommands.
package cli

import (
	"fmt"
	"io"

	"timezone-utility/internal/clock"
	"timezone-utility/internal/data"
	"timezone-utility/internal/output"
)

const version = "0.1.0"

// App bundles dependencies for command execution.
type App struct {
	Store *data.Store
	Clock clock.Clock
	Out   io.Writer
	Err   io.Writer
}

// Run dispatches the command-line arguments and returns the process exit code.
func Run(args []string, app App) int {
	if len(args) == 0 {
		return runNow(nil, app)
	}
	switch args[0] {
	case "now":
		return runNow(args[1:], app)
	case "convert":
		return runConvert(args[1:], app)
	case "meeting":
		return runMeeting(args[1:], app)
	case "overlap":
		return runOverlap(args[1:], app)
	case "--version", "-v", "version":
		fmt.Fprintln(app.Out, version)
		return 0
	case "--help", "-h", "help":
		printUsage(app.Out)
		return 0
	default:
		fmt.Fprintf(app.Err, "error: 未知命令 %q\n\n", args[0])
		printUsage(app.Err)
		return 1
	}
}

func emit(results []data.QueryResult, jsonOut bool, app App) int {
	var err error
	if jsonOut {
		err = output.WriteJSON(app.Out, results)
	} else {
		err = output.WriteText(app.Out, results)
	}
	if err != nil {
		output.WriteError(app.Err, err)
		return 2
	}
	return 0
}

func printUsage(w io.Writer) {
	fmt.Fprint(w, `timezone-utility — 查询任意地点当前时间，支持跨时区换算与会议排程

用法:
  timezone-utility now [query...]                             查询当前时间（时区/地名/邮编）
  timezone-utility convert --from <loc> --to <loc>... <time>  跨时区时间换算
  timezone-utility meeting --from <loc> <time> <loc>...       会议时间展示
  timezone-utility overlap --window <loc>=HH:mm-HH:mm ...     查找重叠可用时段

全局标志:
  --json     输出结构化 JSON
  --version  显示版本
  --help     显示帮助

时间格式: YYYY-MM-DD HH:mm 或 RFC3339
`)
}
