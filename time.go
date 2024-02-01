package nillabletypes

import (
	"bytes"
	"database/sql/driver"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
)

type Time struct {
	v           time.Time
	present     bool
	initialized bool
}

func NewTime(v time.Time) Time {
	return Time{v: v, present: true, initialized: true}
}

func NilTime() Time {
	return Time{present: false, initialized: true}
}

// Value implements the driver.Valuer interface
func (v Time) Value() (driver.Value, error) { //nolint:unparam
	if !v.present {
		return nil, nil
	}
	return v.v, nil
}

func (v Time) Time() time.Time {
	return v.v
}

func (v Time) Nil() bool {
	return !v.present
}

// UnmarshalJSON implements json.Unmarshaler
func (v *Time) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte{'n', 'u', 'l', 'l'}) {
		v.present = false
		v.initialized = true
		return nil
	}
	if err := v.v.UnmarshalJSON(data); err != nil {
		return errors.WithStack(err)
	}
	v.present = true
	v.initialized = true
	return nil
}

// MarshalJSON implements json.Marshaler
func (v Time) MarshalJSON() ([]byte, error) {
	if !v.initialized || !v.present {
		return []byte{'n', 'u', 'l', 'l'}, nil
	}
	return v.v.MarshalJSON()
}

// Scan implements sql.Scanner
func (v *Time) Scan(src any) error {
	switch t := src.(type) {
	case nil:
		*v = Time{present: false, initialized: true}
	case time.Time:
		*v = Time{v: t, present: true, initialized: true}
	case *time.Time:
		if t == nil {
			*v = Time{present: false, initialized: true}
		} else {
			*v = Time{v: *t, present: true, initialized: true}
		}
	default:
		return errors.Errorf("cannot scan value %v to a time", src)
	}

	return nil
}

func (v *Time) ScanTimestamp(src pgtype.Timestamp) error { //nolint:unparam
	*v = Time{v: src.Time, present: src.Valid, initialized: true}
	return nil
}

func (v *Time) ScanTimestamptz(src pgtype.Timestamptz) error { //nolint:unparam
	*v = Time{v: src.Time, present: src.Valid, initialized: true}
	return nil
}

func (v Time) TimestampValue() (pgtype.Timestamp, error) { //nolint:unparam
	return pgtype.Timestamp{Time: v.v, Valid: v.present}, nil
}

func (v Time) TimestamptzValue() (pgtype.Timestamptz, error) { //nolint:unparam
	return pgtype.Timestamptz{Time: v.v, Valid: v.present}, nil
}
