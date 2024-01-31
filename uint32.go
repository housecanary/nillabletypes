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

// Uint32 represents a nil-able int
type Uint32 struct {
	v           uint32
	present     bool
	initialized bool
}

// NewUint32 makes a new non-nil Uint32
func NewUint32(v uint32) Uint32 {
	return Uint32{v: v, present: true, initialized: true}
}

// NilUint32 makes a new nil Uint32
func NilUint32() Uint32 {
	return Uint32{0, false, true}
}

// Uint32 returns the built-in Uint32 value
func (v Uint32) Uint32() uint32 {
	return v.v
}

// Nil returns whether this scalar is nil
func (v Uint32) Nil() bool {
	return !v.present
}

// String implements the fmt.Stringer interface
func (v Uint32) String() string {
	return strconv.FormatUint(uint64(v.v), 10)
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (v *Uint32) UnmarshalJSON(data []byte) error {
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

	i, err := float64ToUint32(t)
	if err != nil {
		return err
	}
	v.v = i
	v.present = true
	v.initialized = true
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (v Uint32) MarshalJSON() ([]byte, error) {
	if !v.initialized || !v.present {
		return []byte{'n', 'u', 'l', 'l'}, nil
	}
	return json.Marshal(v.v)
}

// Value implements the driver.Valuer interface
func (v Uint32) Value() (driver.Value, error) {
	if !v.present {
		return nil, nil
	}
	return v.v, nil
}

// Scan implements the sql.Scanner interface
func (v *Uint32) Scan(src interface{}) error {
	if src == nil {
		*v = Uint32{present: false, initialized: true}
		return nil
	}
	switch t := src.(type) {
	case int64:
		i, err := int64ToUint32(t)
		if err != nil {
			return err
		}
		*v = Uint32{v: i, present: true, initialized: true}
		return nil
	case float64:
		val, err := float64ToUint32(t)
		if err != nil {
			return err
		}
		*v = Uint32{v: val, present: true, initialized: true}
		return nil
	case bool:
		if t {
			*v = Uint32{v: 1, present: true, initialized: true}
		} else {
			*v = Uint32{v: 0, present: true, initialized: true}
		}
		return nil
	case []byte:
		s := string(t)
		return v.scanString(s)
	case string:
		return v.scanString(t)
	}
	return errors.Errorf("cannot scan value %[1]v of type %[1]T to a uint", src)
}

func (v *Uint32) scanString(src string) error {
	f, err := strconv.ParseFloat(src, 32)
	if err != nil {
		return errors.Errorf("cannot scan value %[1]v of type %[1]T to a uint", src)
	}
	i, err := float64ToUint32(f)
	if err != nil {
		return err
	}

	*v = Uint32{v: i, present: true, initialized: true}
	return nil
}

func int64ToUint32(i int64) (uint32, error) {
	if int64(math.MaxUint32) < i || 0 > i {
		return 0, errors.Errorf("value %v outside of the range of Uint32", i)
	}
	return uint32(i), nil
}

func float64ToUint32(f float64) (uint32, error) {
	val := uint32(f)
	if math.Trunc(f) != float64(val) {
		return 0, errors.Errorf("value %f outside of the range of Uint32", f)
	}
	return val, nil
}
