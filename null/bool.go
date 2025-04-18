package null

import (
	"database/sql"
)

// Boolean is a nullable int64.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Boolean struct {
	sql.NullBool
}

// NewBoolean creates a new Boolean
func NewBoolean(b bool) Boolean {
	return Boolean{
		sql.NullBool{
			Bool:  b,
			Valid: true,
		},
	}
}

// NewBooleanFromPtr creates a new Boolean that be null if it is null.
func NewBooleanFromPtr(b *bool) Boolean {
	if b == nil {
		return Boolean{}
	}
	return Boolean{sql.NullBool{Valid: true, Bool: *b}}
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (b Boolean) ValueOrZero() bool {
	if !b.Valid {
		return false
	}
	return b.Bool
}

// Ptr returns a pointer to this Boolean value, or a nil pointer if this Boolean is null.
func (b Boolean) Ptr() *bool {
	if !b.Valid {
		return nil
	}
	return &b.Bool
}

// IsZero returns true for invalid Boolean, for future omitempty support (Go 1.4?)
// A non-null Boolean with a 0 value will not be considered zero.
func (b Boolean) IsZero() bool {
	return !b.Valid
}

// Equal returns true if both integer has the same value or are both null.
func (b Boolean) Equal(other bool) bool {
	return b.Valid && b.Bool == other
}
