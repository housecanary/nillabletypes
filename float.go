package nillabletypes

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
)

// Float represents a nil-able float
type Float struct {
	v           float64
	present     bool
	initialized bool
}

// NewFloat makes a new non-nil Float
func NewFloat(v float64) Float {
	return Float{v: v, present: true, initialized: true}
}

// NilFloat makes a new nil Float
func NilFloat() Float {
	return Float{0, false, true}
}

// Float returns the built-in float64 value
func (v Float) Float() float64 {
	return v.v
}

// Nil returns whether this scalar is nil
func (v Float) Nil() bool {
	return !v.present
}

// String implements the fmt.Stringer interface
func (v Float) String() string {
	if v.Nil() {
		return "nil"
	}
	return strconv.FormatFloat(v.v, 'f', -1, 64)
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (v *Float) UnmarshalJSON(data []byte) error {
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
func (v Float) MarshalJSON() ([]byte, error) {
	if !v.initialized || !v.present {
		return []byte{'n', 'u', 'l', 'l'}, nil
	}
	return json.Marshal(v.v)
}

// Value implements the driver.Valuer interface
func (v Float) Value() (driver.Value, error) {
	if !v.present {
		return nil, nil
	}
	return v.v, nil
}

// Scan implements the sql.Scanner interface
func (v *Float) Scan(src interface{}) error {
	if src == nil {
		*v = Float{present: false, initialized: true}
		return nil
	}
	switch t := src.(type) {
	case float64:
		*v = Float{v: t, present: true, initialized: true}
		return nil
	case int64:
		*v = Float{v: float64(t), present: true, initialized: true}
		return nil
	case bool:
		if t {
			*v = Float{v: 1, present: true, initialized: true}
		} else {
			*v = Float{v: 0, present: true, initialized: true}
		}
		return nil
	case []byte:
		s := string(t)
		return v.scanString(s)
	case string:
		return v.scanString(t)
	}

	return errors.Errorf("cannot scan value %[1]v of type %[1]T to a float", src)
}

func (v *Float) scanString(src string) error {
	f, err := strconv.ParseFloat(src, 64)
	if err != nil {
		return errors.Errorf("cannot scan value %[1]v of type %[1]T to a float", src)
	}

	*v = Float{v: f, present: true, initialized: true}
	return nil
}
