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

// Int32 represents a nil-able int
type Int32 struct {
	v           int32
	present     bool
	initialized bool
}

// NewInt32 makes a new non-nil Int
func NewInt32(v int32) Int32 {
	return Int32{v: v, present: true, initialized: true}
}

// NilInt32 makes a new nil Int
func NilInt32() Int32 {
	return Int32{0, false, true}
}

// Int32 returns the built-in int32 value
func (v Int32) Int32() int32 {
	return v.v
}

// Nil returns whether this scalar is nil
func (v Int32) Nil() bool {
	return !v.present
}

// String implements the fmt.Stringer interface
func (v Int32) String() string {
	return strconv.FormatInt(int64(v.v), 10)
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (v *Int32) UnmarshalJSON(data []byte) error {
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

	i, err := float64ToInt32(t)
	if err != nil {
		return err
	}
	v.v = i
	v.present = true
	v.initialized = true
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (v Int32) MarshalJSON() ([]byte, error) {
	if !v.initialized || !v.present {
		return []byte{'n', 'u', 'l', 'l'}, nil
	}
	return json.Marshal(v.v)
}

// Value implements the driver.Valuer interface
func (v Int32) Value() (driver.Value, error) {
	if !v.present {
		return nil, nil
	}
	return v.v, nil
}

// Scan implements the sql.Scanner interface
func (v *Int32) Scan(src interface{}) error {
	if src == nil {
		*v = Int32{present: false, initialized: true}
		return nil
	}
	switch t := src.(type) {
	case int64:
		i, err := int64ToInt32(t)
		if err != nil {
			return err
		}
		*v = Int32{v: i, present: true, initialized: true}
		return nil
	case float64:
		val, err := float64ToInt32(t)
		if err != nil {
			return err
		}
		*v = Int32{v: val, present: true, initialized: true}
		return nil
	case bool:
		if t {
			*v = Int32{v: 1, present: true, initialized: true}
		} else {
			*v = Int32{v: 0, present: true, initialized: true}
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

func (v *Int32) scanString(src string) error {
	f, err := strconv.ParseFloat(src, 32)
	if err != nil {
		return errors.Errorf("cannot scan value %[1]v of type %[1]T to an int", src)
	}
	i, err := float64ToInt32(f)
	if err != nil {
		return err
	}

	*v = Int32{v: i, present: true, initialized: true}
	return nil
}

func int64ToInt32(i int64) (int32, error) {
	if math.MaxInt32 < i || math.MinInt32 > i {
		return 0, errors.Errorf("value %v outside of the range of int32", i)
	}
	return int32(i), nil
}

func float64ToInt32(f float64) (int32, error) {
	val := int32(f)
	if math.Trunc(f) != float64(val) {
		return 0, errors.Errorf("value %f outside of the range of int32", f)
	}
	return val, nil
}
