package null

import (
	"database/sql"
)

// Int is a nullable int64.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Int struct {
	sql.NullInt64
}

// NewInt creates a new Int
func NewInt(i int64) Int {
	return Int{
		sql.NullInt64{
			Int64: i,
			Valid: true,
		},
	}
}

// NewIntFromPtr creates a new Int that be null if it is null.
func NewIntFromPtr(i *int64) Int {
	if i == nil {
		return Int{}
	}
	return Int{sql.NullInt64{Valid: true, Int64: *i}}
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (i Int) ValueOrZero() int64 {
	if !i.Valid {
		return 0
	}
	return i.Int64
}

// Ptr returns a pointer to this Int value, or a nil pointer if this Int is null.
func (i Int) Ptr() *int64 {
	if !i.Valid {
		return nil
	}
	return &i.Int64
}

// IsZero returns true for invalid Int, for future omitempty support (Go 1.4?)
// A non-null Int with a 0 value will not be considered zero.
func (i Int) IsZero() bool {
	return !i.Valid
}

// Equal returns true if both integer has the same value or are both null.
func (i Int) Equal(other int64) bool {
	return i.Valid && i.Int64 == other
}
