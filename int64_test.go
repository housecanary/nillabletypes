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
	"reflect"
	"testing"
	"time"
)

func TestNewInt64(t *testing.T) {
	tests := []struct {
		name string
		give int64
		want Int64
	}{
		{
			name: "Zero",
			give: 0,
			want: Int64{v: 0, present: true, initialized: true},
		},
		{
			name: "Positive",
			give: 17,
			want: Int64{v: 17, present: true, initialized: true},
		},
		{
			name: "Negative",
			give: -8,
			want: Int64{v: -8, present: true, initialized: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInt64(tt.give); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNilInt64(t *testing.T) {
	tests := []struct {
		name string
		want Int64
	}{
		{"Nil", Int64{v: 0, present: false, initialized: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NilInt64(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NilInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64_Int(t *testing.T) {
	tests := []struct {
		name string
		give Int64
		want int64
	}{
		{
			name: "Non-nil Zero",
			give: Int64{v: 0, present: true, initialized: true},
			want: 0,
		},
		{
			name: "Nil",
			give: Int64{v: 0, present: false, initialized: true},
			want: 0,
		},
		{
			name: "Positive",
			give: Int64{v: 17, present: true, initialized: true},
			want: 17,
		},
		{
			name: "Negative",
			give: Int64{v: -8, present: true, initialized: true},
			want: -8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.give.Int64(); got != tt.want {
				t.Errorf("Int64.Int64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64_Nil(t *testing.T) {
	tests := []struct {
		name string
		give Int64
		want bool
	}{
		{
			name: "Nil",
			give: Int64{present: false, initialized: true},
			want: true,
		},
		{
			name: "Not Nil",
			give: Int64{v: 0, present: true, initialized: true},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.give.Nil(); got != tt.want {
				t.Errorf("Int64.Nil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64_String(t *testing.T) {
	tests := []struct {
		name string
		give Int64
		want string
	}{
		{
			name: "Non-nil Zero",
			give: Int64{v: 0, present: true, initialized: true},
			want: "0",
		},
		{
			name: "Nil",
			give: Int64{v: 0, present: false, initialized: true},
			want: "0",
		},
		{
			name: "Positive",
			give: Int64{v: 98, present: true, initialized: true},
			want: "98",
		},
		{
			name: "Negative",
			give: Int64{v: -1000, present: true, initialized: true},
			want: "-1000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.give.String(); got != tt.want {
				t.Errorf("Int64.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		give    []byte
		want    Int64
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
			want:    Int64{v: 3, present: true, initialized: true},
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
			want:    Int64{v: 100, present: true, initialized: true},
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
			want:    Int64{v: 0, present: false, initialized: true},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &Int64{}
			if err := got.UnmarshalJSON(tt.give); (err != nil) != tt.wantErr {
				t.Errorf("Int64.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("Int64.UnmarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		give    Int64
		want    []byte
		wantErr bool
	}{
		{
			name:    "Present and Initialized",
			give:    Int64{v: 7, present: true, initialized: true},
			want:    toBytes(7),
			wantErr: false,
		},
		{
			name:    "Not Present",
			give:    Int64{v: 7, present: false, initialized: true},
			want:    toBytes(nil),
			wantErr: false,
		},
		{
			name:    "Not Initialized",
			give:    Int64{v: 7, present: true, initialized: false},
			want:    toBytes(nil),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.give.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Int64.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int64.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64_Value(t *testing.T) {
	tests := []struct {
		name    string
		give    Int64
		want    driver.Value
		wantErr bool
	}{
		{
			name:    "Nil",
			give:    Int64{present: false, initialized: true},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "Not Nil",
			give:    Int64{v: 65, present: true, initialized: true},
			want:    int64(65),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.give.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Int64.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int64.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64_Scan(t *testing.T) {
	tests := []struct {
		name    string
		give    interface{}
		want    Int64
		wantErr bool
	}{
		{
			name:    "Nil",
			give:    nil,
			want:    Int64{present: false, initialized: true},
			wantErr: false,
		},
		{
			name:    "Bool (True)",
			give:    true,
			want:    Int64{v: 1, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Bool (False)",
			give:    false,
			want:    Int64{v: 0, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Int",
			give:    int64(123),
			want:    Int64{v: 123, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Float",
			give:    123.456,
			want:    Int64{v: 123, present: true, initialized: true},
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
			want:    Int64{v: 3, present: true, initialized: true},
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
			want:    Int64{v: 92, present: true, initialized: true},
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
			got := &Int64{}
			if err := got.Scan(tt.give); (err != nil) != tt.wantErr {
				t.Errorf("Int64.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("Int64.Scan() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
