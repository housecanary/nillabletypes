package nillabletypes

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"math"
	"strconv"

	"github.com/pkg/errors"
)

// Int represents a nil-able int
type Int struct {
	v           int64
	present     bool
	initialized bool
}

// NewInt makes a new non-nil Int
func NewInt(v int64) Int {
	return Int{v: v, present: true, initialized: true}
}

// NilInt makes a new nil Int
func NilInt() Int {
	return Int{0, false, true}
}

// Int returns the built-in int64 value
func (v Int) Int() int64 {
	return v.v
}

// Nil returns whether this scalar is nil
func (v Int) Nil() bool {
	return !v.present
}

// String implements the fmt.Stringer interface
func (v Int) String() string {
	if v.Nil() {
		return "nil"
	}
	return strconv.FormatInt(v.v, 10)
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (v *Int) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte{'n', 'u', 'l', 'l'}) {
		v.present = false
		v.initialized = true
		return nil
	}
	var t float64
	err := json.Unmarshal(data, &t)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := checkSize(t); err != nil {
		return err
	}
	v.v = int64(t)
	v.present = true
	v.initialized = true
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (v Int) MarshalJSON() ([]byte, error) {
	if !v.initialized || !v.present {
		return []byte{'n', 'u', 'l', 'l'}, nil
	}
	return json.Marshal(v.v)
}

// Value implements the driver.Valuer interface
func (v Int) Value() (driver.Value, error) {
	if !v.present {
		return nil, nil
	}
	return v.v, nil
}

// Scan implements the sql.Scanner interface
func (v *Int) Scan(src interface{}) error {
	if src == nil {
		*v = Int{present: false, initialized: true}
		return nil
	}
	switch t := src.(type) {
	case int64:
		*v = Int{v: t, present: true, initialized: true}
		return nil
	case float64:
		if err := checkSize(t); err != nil {
			return err
		}
		*v = Int{v: int64(t), present: true, initialized: true}
		return nil
	case bool:
		if t {
			*v = Int{v: 1, present: true, initialized: true}
		} else {
			*v = Int{v: 0, present: true, initialized: true}
		}
		return nil
	case []byte:
		s := string(t)
		return v.scanString(s)
	case string:
		return v.scanString(t)
	}
	return errors.Errorf("cannot scan value %[1]v of type %[1]T to an int", src)
}

func (v *Int) scanString(src string) error {
	f, err := strconv.ParseFloat(src, 64)
	if err != nil {
		return errors.Errorf("cannot scan value %[1]v of type %[1]T to an int", src)
	}

	if err := checkSize(f); err != nil {
		return err
	}

	*v = Int{v: int64(f), present: true, initialized: true}
	return nil
}

func checkSize(f float64) error {
	if math.Trunc(f) != float64(int64(f)) {
		return errors.Errorf("value %f outside of the range of int64", f)
	}
	return nil
}
