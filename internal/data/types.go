// Package data defines the domain types and the embedded location datasets.
package data

// Location represents a resolvable geographic location.
type Location struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Aliases     []string `json:"aliases,omitempty"`
	Country     string   `json:"country"`
	Timezone    string   `json:"timezone"`
	PostalCodes []string `json:"postal_codes,omitempty"`
}

// QueryResult is the output of a single lookup.
type QueryResult struct {
	Query            string    `json:"query"`
	ResolvedLocation *Location `json:"resolved,omitempty"`
	LocalDate        string    `json:"local_date"`
	LocalTime        string    `json:"local_time"`
	Timezone         string    `json:"timezone"`
	UTCOffset        string    `json:"utc_offset"`
	DST              bool      `json:"dst"`
}
