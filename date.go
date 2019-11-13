package nillabletypes

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"regexp"
	"time"

	"github.com/pkg/errors"
)

// TODO(faulkner): technically this is incorrect since it'll match "2000-12-3456789"; useful if we want to trim the time from a datetime, but if we don't have that usecase then perhaps this should be more strict.
var datePattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}`)

// Date represents a nil-able date encoded as ISO
type Date struct {
	v           string
	present     bool
	initialized bool
}

// NewDate makes a new non-nil Date
func NewDate(v string) Date {
	return Date{v: v, present: true, initialized: true}
}

// NilDate makes a new nil Date
func NilDate() Date {
	return Date{v: "", present: false, initialized: true}
}

// Nil returns whether this scalar is nil
func (v Date) Nil() bool {
	return !v.present
}

// DaysAgo returns the number of days elapsed since Date
func (v Date) DaysAgo() (Int, error) {
	if v.Nil() {
		return Int{}, nil
	}

	t, err := time.Parse("2006-01-02", v.String())
	if err != nil {
		return Int{}, err
	}

	return NewInt(int64(time.Since(t).Hours() / 24)), nil
}

// String implements the fmt.Stringer interface
func (v Date) String() string {
	if v.Nil() {
		return "nil"
	}
	return v.v
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (v *Date) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte{'n', 'u', 'l', 'l'}) {
		*v = Date{present: false, initialized: true}
		return nil
	}

	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return errors.WithStack(err)
	}

	if f := datePattern.FindString(s); f == "" {
		return errors.Errorf("value %v is not a valid date", s)
	}

	*v = Date{v: s, present: true, initialized: true}
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (v Date) MarshalJSON() ([]byte, error) {
	if !v.initialized || !v.present {
		return []byte{'n', 'u', 'l', 'l'}, nil
	}
	return json.Marshal(v.v)
}

// Value implements the driver.Valuer interface
func (v Date) Value() (driver.Value, error) {
	if !v.present {
		return nil, nil
	}
	return v.v, nil
}

// Scan implements the sql.Scanner interface
func (v *Date) Scan(src interface{}) error {
	if src == nil {
		*v = Date{present: false, initialized: true}
		return nil
	}
	if t, ok := src.(time.Time); ok {
		*v = Date{v: t.Format("2006-01-02"), present: true, initialized: true}
		return nil
	}
	return errors.Errorf("cannot scan value %v to a date", src)
}
