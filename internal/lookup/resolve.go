// Package lookup resolves queries to timezones and computes current times.
package lookup

import (
	"fmt"
	"strings"
	"time"

	"timezone-utility/internal/data"
)

// Resolve resolves a raw query to an IANA timezone name and optional location.
// It accepts a timezone identifier, a place name, or a postal code.
func Resolve(store *data.Store, query string) (string, *data.Location, error) {
	q := strings.TrimSpace(query)
	if q == "" {
		return "", nil, fmt.Errorf("空查询")
	}
	if _, err := time.LoadLocation(q); err == nil {
		return q, nil, nil
	}
	if isNumeric(q) {
		loc, err := store.ResolvePostal(q)
		if err != nil {
			return "", nil, err
		}
		return loc.Timezone, loc, nil
	}
	loc, err := store.ResolvePlace(q)
	if err != nil {
		return "", nil, fmt.Errorf("无法识别查询 %q：不是有效时区、地点名称或邮政编码", q)
	}
	return loc.Timezone, loc, nil
}

func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
