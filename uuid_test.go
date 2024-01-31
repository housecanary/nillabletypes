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
	"database/sql/driver"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/encoding/json"
	"github.com/stretchr/testify/assert"
)

const stubUUIDString = "11111111-1111-1111-1111-111111111111"

var stubUUID = uuid.MustParse(stubUUIDString)

func TestNewUUID(t *testing.T) {
	assert.Equal(t, UUID{v: uuid.UUID{}, present: true, initialized: true}, NewUUID(uuid.UUID{}))
	assert.Equal(t, UUID{v: stubUUID, present: true, initialized: true}, NewUUID(stubUUID))
}

func TestNilUUID(t *testing.T) {
	assert.Equal(t, UUID{present: false, initialized: true}, NilUUID())
}

func TestUUID_Nil(t *testing.T) {
	assert.True(t, UUID{present: false, initialized: true}.Nil())
	assert.False(t, UUID{v: stubUUID, present: true, initialized: true}.Nil())
}

func TestUUID_UUID(t *testing.T) {
	tests := []struct {
		name string
		give UUID
		want uuid.UUID
	}{
		{
			name: "Non-nil Empty",
			give: UUID{v: uuid.UUID{}, present: true, initialized: true},
			want: uuid.UUID{},
		},
		{
			name: "Nil",
			give: UUID{present: false, initialized: true},
			want: uuid.UUID{},
		},
		{
			name: "Not Empty",
			give: UUID{v: stubUUID, present: true, initialized: true},
			want: stubUUID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.give.UUID())
		})
	}
}

func TestUUID_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		give    []byte
		want    UUID
		wantErr bool
	}{
		{
			name:    "UUID",
			give:    toJSONBytes(stubUUIDString),
			want:    UUID{v: stubUUID, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Floating Point Number",
			give:    toJSONBytes(3.14159),
			wantErr: true,
		},
		{
			name:    "Integer",
			give:    toJSONBytes(100),
			wantErr: true,
		},
		{
			name:    "Boolean",
			give:    toJSONBytes(true),
			wantErr: true,
		},
		{
			name:    "Null",
			give:    toJSONBytes(nil),
			want:    UUID{present: false, initialized: true},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &UUID{}
			err := got.UnmarshalJSON(tt.give)
			assertWantError(t, tt.wantErr, err)
			assert.Equal(t, tt.want, *got)
		})
	}
}

func TestUUID_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		give    UUID
		want    []byte
		wantErr bool
	}{
		{
			name:    "Present and Initialized",
			give:    UUID{v: stubUUID, present: true, initialized: true},
			want:    toJSONBytes(stubUUIDString),
			wantErr: false,
		},
		{
			name:    "Not Present",
			give:    UUID{v: stubUUID, present: false, initialized: true},
			want:    toJSONBytes(nil),
			wantErr: false,
		},
		{
			name:    "Not Initialized",
			give:    UUID{v: stubUUID, present: true, initialized: false},
			want:    toJSONBytes(nil),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.give.MarshalJSON()
			assertWantError(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUUID_Value(t *testing.T) {
	tests := []struct {
		name    string
		give    UUID
		want    driver.Value
		wantErr bool
	}{
		{
			name:    "Nil",
			give:    UUID{present: false, initialized: true},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "Not Nil",
			give:    UUID{v: stubUUID, present: true, initialized: true},
			want:    stubUUID,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.give.Value()
			assertWantError(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUUID_Scan(t *testing.T) {
	tests := []struct {
		name    string
		give    any
		want    UUID
		wantErr bool
	}{
		{
			name:    "Nil",
			give:    nil,
			want:    UUID{present: false, initialized: true},
			wantErr: false,
		},
		{
			name:    "Bool",
			give:    true,
			wantErr: true,
		},
		{
			name:    "Int",
			give:    int64(123),
			wantErr: true,
		},
		{
			name:    "Float",
			give:    123.456,
			wantErr: true,
		},
		{
			name:    "Byte Slice",
			give:    stubUUID[:],
			want:    UUID{v: stubUUID, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "UUID",
			give:    stubUUIDString,
			want:    UUID{v: stubUUID, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Time",
			give:    time.Now(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &UUID{}
			err := got.Scan(tt.give)
			assertWantError(t, tt.wantErr, err)
			assert.Equal(t, tt.want, *got)
		})
	}
}

func assertWantError(t *testing.T, wantErr bool, err error) {
	t.Helper()
	if wantErr {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
	}
}

func toJSONBytes(val any) []byte {
	bytes, err := json.Marshal(val)
	if err != nil {
		panic(fmt.Sprintf("Failed to convert value (%v) of type %T to bytes: %v\n", val, val, err))
	}
	return bytes
}
