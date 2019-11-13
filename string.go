package nillabletypes

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/pkg/errors"
)

// String represents a nil-able string
type String struct {
	v           string
	present     bool
	initialized bool
}

// NewString makes a new non-nil String
func NewString(v string) String {
	return String{v: v, present: true, initialized: true}
}

// NilString makes a new nil String
func NilString() String {
	return String{"", false, true}
}

// Nil returns whether this scalar is nil
func (v String) Nil() bool {
	return !v.present
}

// String implements the fmt.Stringer interface
func (v String) String() string {
	return v.v
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (v *String) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte{'n', 'u', 'l', 'l'}) {
		v.present = false
		v.initialized = true
		return nil
	}
	err := json.Unmarshal(data, &v.v)
	if err != nil {
		return errors.WithStack(err)
	}
	v.present = true
	v.initialized = true
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (v String) MarshalJSON() ([]byte, error) {
	if !v.initialized || !v.present {
		return []byte{'n', 'u', 'l', 'l'}, nil
	}
	return json.Marshal(v.v)
}

// Value implements the driver.Valuer interface
func (v String) Value() (driver.Value, error) {
	if !v.present {
		return nil, nil
	}
	return v.v, nil
}

// Scan implements the sql.Scanner interface
func (v *String) Scan(src interface{}) error {
	if src == nil {
		*v = String{present: false, initialized: true}
		return nil
	}
	switch t := src.(type) {
	case []byte:
		*v = String{present: true, v: string(t), initialized: true}
		return nil
	case string:
		*v = String{present: true, v: t, initialized: true}
		return nil
	}
	return errors.Errorf("cannot scan value %v to a string", src)
}
