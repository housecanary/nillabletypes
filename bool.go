package nillabletypes

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
)

// Bool represents a nil-able bool
type Bool struct {
	v           bool
	present     bool
	initialized bool
}

// NewBool makes a new non-nil Bool
func NewBool(v bool) Bool {
	return Bool{v: v, present: true, initialized: true}
}

// NilBool makes a new nil Bool
func NilBool() Bool {
	return Bool{v: false, present: false, initialized: true}
}

// Bool returns the built-in bool value
func (v Bool) Bool() bool {
	return v.v
}

// Nil returns whether this scalar is nil
func (v Bool) Nil() bool {
	return !v.present
}

// String implements the fmt.Stringer interface
func (v Bool) String() string {
	if !v.initialized || !v.present {
		return "nil"
	}
	return strconv.FormatBool(v.v)
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (v *Bool) UnmarshalJSON(data []byte) error {
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
func (v Bool) MarshalJSON() ([]byte, error) {
	if !v.initialized || !v.present {
		return []byte{'n', 'u', 'l', 'l'}, nil
	}
	return json.Marshal(v.v)
}

// Value implements the driver.Valuer interface
func (v Bool) Value() (driver.Value, error) {
	if !v.present {
		return nil, nil
	}
	return v.v, nil
}

// Scan implements the sql.Scanner interface
func (v *Bool) Scan(src interface{}) error {
	if src == nil {
		*v = Bool{present: false, initialized: true}
		return nil
	}
	switch t := src.(type) {
	case bool:
		*v = Bool{v: t, present: true, initialized: true}
		return nil
	case int64:
		if t == 0 {
			*v = Bool{v: false, present: true, initialized: true}
		} else {
			*v = Bool{v: true, present: true, initialized: true}
		}
		return nil
	case float64:
		if t == 0 {
			*v = Bool{v: false, present: true, initialized: true}
		} else {
			*v = Bool{v: true, present: true, initialized: true}
		}
		return nil
	}
	return errors.Errorf("cannot scan value %[1]v of type %[1]T to a bool", src)
}
