package null

import (
	"database/sql"
	"strings"
)

// String is a nullable string. It supports SQL and JSON serialization.
// It will marshal to null if null. Blank string input will be considered null.
type String struct {
	sql.NullString
}

// NewStringFromPtr creates a new String that be null if s is nil.
func NewStringFromPtr(s *string) String {
	if s == nil {
		return String{}
	}
	
	return String{sql.NullString{Valid: true, String: *s}}
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (s String) ValueOrZero() string {
	if !s.Valid {
		return ""
	}
	return s.String
}

// NewString creates a new String
func NewString(s string) String {
	if s == "" {
		return String{}
	}
	
	return String{
		sql.NullString{
			String: s,
			Valid:  true,
		},
	}
}

// Ptr returns a pointer to this String's value, or a nil pointer if this String is null.
func (s String) Ptr() *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}

// IsZero returns true for null strings, for potential future omitempty support.
func (s String) IsZero() bool {
	return !s.Valid
}

// Equal returns true if both strings have the same value or are both null.
func (s String) Equal(other string) bool {
	return s.Valid && strings.EqualFold(s.String, other)
}
