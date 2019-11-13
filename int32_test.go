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

func TestNewInt32(t *testing.T) {
	tests := []struct {
		name string
		give int32
		want Int32
	}{
		{
			name: "Zero",
			give: 0,
			want: Int32{v: 0, present: true, initialized: true},
		},
		{
			name: "Positive",
			give: 17,
			want: Int32{v: 17, present: true, initialized: true},
		},
		{
			name: "Negative",
			give: -8,
			want: Int32{v: -8, present: true, initialized: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInt32(tt.give); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNilInt32(t *testing.T) {
	tests := []struct {
		name string
		want Int32
	}{
		{"Nil", Int32{v: 0, present: false, initialized: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NilInt32(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NilInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt32_Int(t *testing.T) {
	tests := []struct {
		name string
		give Int32
		want int32
	}{
		{
			name: "Non-nil Zero",
			give: Int32{v: 0, present: true, initialized: true},
			want: 0,
		},
		{
			name: "Nil",
			give: Int32{v: 0, present: false, initialized: true},
			want: 0,
		},
		{
			name: "Positive",
			give: Int32{v: 17, present: true, initialized: true},
			want: 17,
		},
		{
			name: "Negative",
			give: Int32{v: -8, present: true, initialized: true},
			want: -8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.give.Int32(); got != tt.want {
				t.Errorf("Int32.Int32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt32_Nil(t *testing.T) {
	tests := []struct {
		name string
		give Int32
		want bool
	}{
		{
			name: "Nil",
			give: Int32{present: false, initialized: true},
			want: true,
		},
		{
			name: "Not Nil",
			give: Int32{v: 0, present: true, initialized: true},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.give.Nil(); got != tt.want {
				t.Errorf("Int32.Nil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt32_String(t *testing.T) {
	tests := []struct {
		name string
		give Int32
		want string
	}{
		{
			name: "Non-nil Zero",
			give: Int32{v: 0, present: true, initialized: true},
			want: "0",
		},
		{
			name: "Nil",
			give: Int32{v: 0, present: false, initialized: true},
			want: "0",
		},
		{
			name: "Positive",
			give: Int32{v: 98, present: true, initialized: true},
			want: "98",
		},
		{
			name: "Negative",
			give: Int32{v: -1000, present: true, initialized: true},
			want: "-1000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.give.String(); got != tt.want {
				t.Errorf("Int32.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt32_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		give    []byte
		want    Int32
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
			want:    Int32{v: 3, present: true, initialized: true},
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
			want:    Int32{v: 100, present: true, initialized: true},
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
			want:    Int32{v: 0, present: false, initialized: true},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &Int32{}
			if err := got.UnmarshalJSON(tt.give); (err != nil) != tt.wantErr {
				t.Errorf("Int32.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("Int32.UnmarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt32_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		give    Int32
		want    []byte
		wantErr bool
	}{
		{
			name:    "Present and Initialized",
			give:    Int32{v: 7, present: true, initialized: true},
			want:    toBytes(7),
			wantErr: false,
		},
		{
			name:    "Not Present",
			give:    Int32{v: 7, present: false, initialized: true},
			want:    toBytes(nil),
			wantErr: false,
		},
		{
			name:    "Not Initialized",
			give:    Int32{v: 7, present: true, initialized: false},
			want:    toBytes(nil),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.give.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Int32.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int32.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt32_Value(t *testing.T) {
	tests := []struct {
		name    string
		give    Int32
		want    driver.Value
		wantErr bool
	}{
		{
			name:    "Nil",
			give:    Int32{present: false, initialized: true},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "Not Nil",
			give:    Int32{v: 65, present: true, initialized: true},
			want:    int32(65),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.give.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Int32.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int32.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt32_Scan(t *testing.T) {
	tests := []struct {
		name    string
		give    interface{}
		want    Int32
		wantErr bool
	}{
		{
			name:    "Nil",
			give:    nil,
			want:    Int32{present: false, initialized: true},
			wantErr: false,
		},
		{
			name:    "Bool (True)",
			give:    true,
			want:    Int32{v: 1, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Bool (False)",
			give:    false,
			want:    Int32{v: 0, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Int (Max)",
			give:    int64(2147483647),
			want:    Int32{v: 2147483647, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Int (Too Large)",
			give:    int64(2147483648),
			wantErr: true,
		},
		{
			name:    "Int (Min)",
			give:    int64(-2147483648),
			want:    Int32{v: -2147483648, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Int (Too Small)",
			give:    int64(-2147483649),
			wantErr: true,
		},
		{
			name:    "Float",
			give:    123.456,
			want:    Int32{v: 123, present: true, initialized: true},
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
			want:    Int32{v: 3, present: true, initialized: true},
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
			want:    Int32{v: 92, present: true, initialized: true},
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
			got := &Int32{}
			if err := got.Scan(tt.give); (err != nil) != tt.wantErr {
				t.Errorf("Int32.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("Int32.Scan() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
