package null

import (
	"database/sql"
	"time"
)

// Time is a nullable float64.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Time struct {
	sql.NullTime
}

// NewTime creates a new Time
func NewTime(t time.Time) Time {
	if t.IsZero() {
		return Time{}
	}

	return Time{
		sql.NullTime{
			Time:  t,
			Valid: true,
		},
	}
}

// NewTimeFromPtr creates a new Time that be null if f is nil.
func NewTimeFromPtr(t *time.Time) Time {
	if t == nil {
		return Time{}
	}

	if t.IsZero() {
		return Time{}
	}

	return Time{sql.NullTime{Valid: true, Time: *t}}
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (t Time) ValueOrZero() time.Time {
	if !t.Valid {
		return time.Time{}
	}

	return t.Time
}

// Ptr returns a pointer to this Time's value, or a nil pointer if this Time is null.
func (t Time) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}

	return &t.Time
}

// IsZero returns true for invalid Times, for future omitempty support (Go 1.4?)
// A non-null Time with a 0 value will not be considered zero.
func (t Time) IsZero() bool {
	return !t.Valid || t.Time.IsZero()
}

// Equal returns true if both floats have the same value or are both null.
// Warning: calculations using floating point numbers can result in different ways
// the numbers are stored in memory. Therefore, this function is not suitable to
// compare the result of a calculation. Use this method only to check if the value
// has changed in comparison to some previous value.
func (t Time) Equal(ti time.Time) bool {
	return t.Valid && t.Time.Equal(ti)
}
