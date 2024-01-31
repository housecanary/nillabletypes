// Copyright 2019 HouseCanary, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package nillabletypes

import (
	"bytes"
	"database/sql/driver"
	"math"
	"strconv"

	"github.com/pkg/errors"
	"github.com/segmentio/encoding/json"
)

// Int64 represents a nil-able int
type Int64 struct {
	v           int64
	present     bool
	initialized bool
}

// NewInt64 makes a new non-nil Int
func NewInt64(v int64) Int64 {
	return Int64{v: v, present: true, initialized: true}
}

// NilInt64 makes a new nil Int
func NilInt64() Int64 {
	return Int64{0, false, true}
}

// Int64 returns the built-in int64 value
func (v Int64) Int64() int64 {
	return v.v
}

// Nil returns whether this scalar is nil
func (v Int64) Nil() bool {
	return !v.present
}

// String implements the fmt.Stringer interface
func (v Int64) String() string {
	return strconv.FormatInt(v.v, 10)
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (v *Int64) UnmarshalJSON(data []byte) error {
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

	i, err := float64ToInt64(t)
	if err != nil {
		return err
	}
	v.v = i
	v.present = true
	v.initialized = true
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (v Int64) MarshalJSON() ([]byte, error) {
	if !v.initialized || !v.present {
		return []byte{'n', 'u', 'l', 'l'}, nil
	}
	return json.Marshal(v.v)
}

// Value implements the driver.Valuer interface
func (v Int64) Value() (driver.Value, error) {
	if !v.present {
		return nil, nil
	}
	return v.v, nil
}

// Scan implements the sql.Scanner interface
func (v *Int64) Scan(src interface{}) error {
	if src == nil {
		*v = Int64{present: false, initialized: true}
		return nil
	}
	switch t := src.(type) {
	case int64:
		*v = Int64{v: t, present: true, initialized: true}
		return nil
	case float64:
		i, err := float64ToInt64(t)
		if err != nil {
			return err
		}
		*v = Int64{v: i, present: true, initialized: true}
		return nil
	case bool:
		if t {
			*v = Int64{v: 1, present: true, initialized: true}
		} else {
			*v = Int64{v: 0, present: true, initialized: true}
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

func (v *Int64) scanString(src string) error {
	f, err := strconv.ParseFloat(src, 64)
	if err != nil {
		return errors.Errorf("cannot scan value %[1]v of type %[1]T to an int", src)
	}

	i, err := float64ToInt64(f)
	if err != nil {
		return err
	}

	*v = Int64{v: i, present: true, initialized: true}
	return nil
}

func float64ToInt64(f float64) (int64, error) {
	val := int64(f)
	if math.Trunc(f) != float64(val) {
		return 0, errors.Errorf("value %f outside of the range of int64", f)
	}
	return val, nil
}
