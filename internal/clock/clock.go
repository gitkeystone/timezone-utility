// Package clock provides an injectable time source for deterministic testing.
package clock

import "time"

// Clock returns the current time.
type Clock interface {
	Now() time.Time
}

// RealClock uses the system clock.
type RealClock struct{}

// Now returns the current system time.
func (RealClock) Now() time.Time { return time.Now() }

// FixedClock returns a fixed time, for tests and reproducible output.
type FixedClock struct {
	T time.Time
}

// Now returns the fixed time.
func (f FixedClock) Now() time.Time { return f.T }
