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
	"database/sql/driver"
	"math"
	"reflect"
	"testing"
	"time"
)

func TestNewUint32(t *testing.T) {
	tests := []struct {
		name string
		give uint32
		want Uint32
	}{
		{
			name: "Zero",
			give: 0,
			want: Uint32{v: 0, present: true, initialized: true},
		},
		{
			name: "Positive",
			give: 17,
			want: Uint32{v: 17, present: true, initialized: true},
		},
		{
			name: "Max",
			give: math.MaxUint32,
			want: Uint32{v: math.MaxUint32, present: true, initialized: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUint32(tt.give); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUint32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNilUint32(t *testing.T) {
	tests := []struct {
		name string
		want Uint32
	}{
		{"Nil", Uint32{v: 0, present: false, initialized: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NilUint32(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NilUint32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint32_Uint32(t *testing.T) {
	tests := []struct {
		name string
		give Uint32
		want uint32
	}{
		{
			name: "Non-nil Zero",
			give: Uint32{v: 0, present: true, initialized: true},
			want: 0,
		},
		{
			name: "Nil",
			give: Uint32{v: 0, present: false, initialized: true},
			want: 0,
		},
		{
			name: "Positive",
			give: Uint32{v: 17, present: true, initialized: true},
			want: 17,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.give.Uint32(); got != tt.want {
				t.Errorf("Uint32.Uint32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUInt32_Nil(t *testing.T) {
	tests := []struct {
		name string
		give Uint32
		want bool
	}{
		{
			name: "Nil",
			give: Uint32{present: false, initialized: true},
			want: true,
		},
		{
			name: "Not Nil",
			give: Uint32{v: 0, present: true, initialized: true},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.give.Nil(); got != tt.want {
				t.Errorf("Uint32.Nil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint32_String(t *testing.T) {
	tests := []struct {
		name string
		give Uint32
		want string
	}{
		{
			name: "Non-nil Zero",
			give: Uint32{v: 0, present: true, initialized: true},
			want: "0",
		},
		{
			name: "Nil",
			give: Uint32{v: 0, present: false, initialized: true},
			want: "0",
		},
		{
			name: "Positive",
			give: Uint32{v: 98, present: true, initialized: true},
			want: "98",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.give.String(); got != tt.want {
				t.Errorf("Uint32.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint32_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		give    []byte
		want    Uint32
		wantErr bool
	}{
		{
			name:    "String",
			give:    toBytes("3"),
			wantErr: true,
		},
		{
			name:    "Floating Point Number",
			give:    toBytes(3.14159),
			want:    Uint32{v: 3, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Floating Point Number (Max)",
			give:    toBytes(float64(math.MaxUint32)),
			want:    Uint32{v: math.MaxUint32, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Floating Point Number (Out of Range)",
			give:    toBytes(9999999999999999999.0),
			wantErr: true,
		},
		{
			name:    "Integer",
			give:    toBytes(100),
			want:    Uint32{v: 100, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Boolean",
			give:    toBytes(true),
			wantErr: true,
		},
		{
			name:    "Null",
			give:    toBytes(nil),
			want:    Uint32{v: 0, present: false, initialized: true},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &Uint32{}
			if err := got.UnmarshalJSON(tt.give); (err != nil) != tt.wantErr {
				t.Errorf("Uint32.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("Uint32.UnmarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint32_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		give    Uint32
		want    []byte
		wantErr bool
	}{
		{
			name:    "Present and Initialized",
			give:    Uint32{v: 7, present: true, initialized: true},
			want:    toBytes(7),
			wantErr: false,
		},
		{
			name:    "Not Present",
			give:    Uint32{v: 7, present: false, initialized: true},
			want:    toBytes(nil),
			wantErr: false,
		},
		{
			name:    "Not Initialized",
			give:    Uint32{v: 7, present: true, initialized: false},
			want:    toBytes(nil),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.give.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Uint32.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Uint32.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint32_Value(t *testing.T) {
	tests := []struct {
		name    string
		give    Uint32
		want    driver.Value
		wantErr bool
	}{
		{
			name:    "Nil",
			give:    Uint32{present: false, initialized: true},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "Not Nil",
			give:    Uint32{v: 65, present: true, initialized: true},
			want:    uint32(65),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.give.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Uint32.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Uint32.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint32_Scan(t *testing.T) {
	tests := []struct {
		name    string
		give    interface{}
		want    Uint32
		wantErr bool
	}{
		{
			name:    "Nil",
			give:    nil,
			want:    Uint32{present: false, initialized: true},
			wantErr: false,
		},
		{
			name:    "Bool (True)",
			give:    true,
			want:    Uint32{v: 1, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Bool (False)",
			give:    false,
			want:    Uint32{v: 0, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Int (Max)",
			give:    int64(math.MaxUint32),
			want:    Uint32{v: math.MaxUint32, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Int (Too Large)",
			give:    int64(math.MaxUint32 + 1),
			wantErr: true,
		},
		{
			name:    "Int (Min)",
			give:    int64(0),
			want:    Uint32{v: 0, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Int (Too Small)",
			give:    int64(-1),
			wantErr: true,
		},
		{
			name:    "Float",
			give:    123.456,
			want:    Uint32{v: 123, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Float (Out of Range)",
			give:    9999999999999999999.0,
			wantErr: true,
		},
		{
			name:    "Byte Slice (Valid)",
			give:    []byte{'3', '.', '1', '4'},
			want:    Uint32{v: 3, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Byte Slice (Invalid)",
			give:    []byte{'t', 'o', 'r', 'n', 'a', 'd', 'o'},
			wantErr: true,
		},
		{
			name:    "String (Valid)",
			give:    "92.96",
			want:    Uint32{v: 92, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "String (Out of Range)",
			give:    "9999999999999999999.0",
			wantErr: true,
		},
		{
			name:    "String (Invalid)",
			give:    "twister",
			wantErr: true,
		},
		{
			name:    "Time",
			give:    time.Now(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &Uint32{}
			if err := got.Scan(tt.give); (err != nil) != tt.wantErr {
				t.Errorf("Uint32.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("Uint32.Scan() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
