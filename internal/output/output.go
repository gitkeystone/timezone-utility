// Package output renders query results in a consistent text or JSON format.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"timezone-utility/internal/data"
)

// Results wraps a list of query results for JSON output.
type Results struct {
	Results []data.QueryResult `json:"results"`
}

// WriteText renders human-readable tabular output.
func WriteText(w io.Writer, results []data.QueryResult) error {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "LOCATION\tDATE\tTIME\tUTC OFFSET\tDST\tTIMEZONE")
	for _, r := range results {
		name := r.Timezone
		if r.ResolvedLocation != nil {
			name = r.ResolvedLocation.Name
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%v\t%s\n", name, r.LocalDate, r.LocalTime, r.UTCOffset, r.DST, r.Timezone)
	}
	return tw.Flush()
}

// WriteJSON renders structured JSON output.
func WriteJSON(w io.Writer, results []data.QueryResult) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(Results{Results: results})
}

// WriteError renders a unified error message.
func WriteError(w io.Writer, err error) {
	fmt.Fprintf(w, "error: %s\n", err)
}
