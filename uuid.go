// Copyright 2024 HouseCanary, Inc.
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

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
	"github.com/segmentio/encoding/json"
)

// UUID represents a nil-able UUID
type UUID struct {
	v           uuid.UUID
	present     bool
	initialized bool
}

// NewUUID makes a new non-nil UUID
func NewUUID(v uuid.UUID) UUID {
	return UUID{v: v, present: true, initialized: true}
}

// NilUUID makes a new nil UUID
func NilUUID() UUID {
	return UUID{uuid.Nil, false, true}
}

// Nil returns whether this scalar is nil
func (v UUID) Nil() bool {
	return !v.present
}

// String implements the fmt.Stringer interface
func (v UUID) String() string {
	if !v.present || v.v == uuid.Nil {
		return ""
	}
	return v.v.String()
}

func (v UUID) UUID() uuid.UUID {
	return v.v
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (v *UUID) UnmarshalJSON(data []byte) error {
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
func (v UUID) MarshalJSON() ([]byte, error) {
	if !v.initialized || !v.present {
		return []byte{'n', 'u', 'l', 'l'}, nil
	}
	return json.Marshal(v.v)
}

// Value implements the driver.Valuer interface
func (v UUID) Value() (driver.Value, error) { //nolint:unparam
	if !v.present {
		return nil, nil
	}
	return v.v, nil
}

// Scan implements the sql.Scanner interface
func (v *UUID) Scan(src any) error {
	if src == nil {
		*v = UUID{present: false, initialized: true}
		return nil
	}
	switch t := src.(type) {
	case []byte:
		u, err := uuid.FromBytes(t)
		if err != nil {
			return err
		}
		*v = UUID{present: true, v: u, initialized: true}
		return nil
	case string:
		u, err := uuid.Parse(t)
		if err != nil {
			return err
		}
		*v = UUID{present: true, v: u, initialized: true}
		return nil
	}
	return errors.Errorf("cannot scan value %v to a UUID", src)
}

// ScanUUID implements the pgtype.UUIDScanner interface
func (v *UUID) ScanUUID(src pgtype.UUID) error { //nolint:unparam
	*v = UUID{present: src.Valid, v: src.Bytes, initialized: true}
	return nil
}

// UUIDValue implements the pgtype.UUIDValuer interface
func (v *UUID) UUIDValue() (pgtype.UUID, error) { //nolint:unparam
	return pgtype.UUID{Bytes: v.v, Valid: v.present}, nil
}

func StringsToUUIDs(ids []string) []uuid.UUID {
	v := make([]uuid.UUID, len(ids))
	for i := range ids {
		v[i] = uuid.MustParse(ids[i])
	}
	return v
}
