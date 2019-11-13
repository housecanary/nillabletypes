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

func TestNewString(t *testing.T) {
	tests := []struct {
		name string
		give string
		want String
	}{
		{
			name: "Empty",
			give: "",
			want: String{v: "", present: true, initialized: true},
		},
		{
			name: "Not Empty",
			give: "unicorn",
			want: String{v: "unicorn", present: true, initialized: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewString(tt.give); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNilString(t *testing.T) {
	tests := []struct {
		name string
		want String
	}{
		{"Nil", String{v: "", present: false, initialized: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NilString(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NilString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_Nil(t *testing.T) {
	tests := []struct {
		name string
		give String
		want bool
	}{
		{
			name: "Nil",
			give: String{present: false, initialized: true},
			want: true,
		},
		{
			name: "Not Nil",
			give: String{v: "wooly mammoth", present: true, initialized: true},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.give.Nil(); got != tt.want {
				t.Errorf("String.Nil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_String(t *testing.T) {
	tests := []struct {
		name string
		give String
		want string
	}{
		{
			name: "Non-nil Empty",
			give: String{v: "", present: true, initialized: true},
			want: "",
		},
		{
			name: "Nil",
			give: String{v: "", present: false, initialized: true},
			want: "",
		},
		{
			name: "Not Empty",
			give: String{v: "koala", present: true, initialized: true},
			want: "koala",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.give.String(); got != tt.want {
				t.Errorf("String.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		give    []byte
		want    String
		wantErr bool
	}{
		{
			name:    "String",
			give:    toBytes("habanero"),
			want:    String{v: "habanero", present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Floating Point Number",
			give:    toBytes(3.14159),
			wantErr: true,
		},
		{
			name:    "Integer",
			give:    toBytes(100),
			wantErr: true,
		},
		{
			name:    "Boolean",
			give:    toBytes(true),
			wantErr: true,
		},
		{
			name:    "Null",
			give:    toBytes(nil),
			want:    String{v: "", present: false, initialized: true},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &String{}
			if err := got.UnmarshalJSON(tt.give); (err != nil) != tt.wantErr {
				t.Errorf("String.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("String.UnmarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		give    String
		want    []byte
		wantErr bool
	}{
		{
			name:    "Present and Initialized",
			give:    String{v: "sash", present: true, initialized: true},
			want:    toBytes("sash"),
			wantErr: false,
		},
		{
			name:    "Not Present",
			give:    String{v: "fedora", present: false, initialized: true},
			want:    toBytes(nil),
			wantErr: false,
		},
		{
			name:    "Not Initialized",
			give:    String{v: "monacle", present: true, initialized: false},
			want:    toBytes(nil),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.give.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("String.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("String.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_Value(t *testing.T) {
	tests := []struct {
		name    string
		give    String
		want    driver.Value
		wantErr bool
	}{
		{
			name:    "Nil",
			give:    String{present: false, initialized: true},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "Not Nil",
			give:    String{v: "habidashery", present: true, initialized: true},
			want:    "habidashery",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.give.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("String.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("String.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_Scan(t *testing.T) {
	tests := []struct {
		name    string
		give    interface{}
		want    String
		wantErr bool
	}{
		{
			name:    "Nil",
			give:    nil,
			want:    String{present: false, initialized: true},
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
			give:    []byte{'h', 'i', 'p', 'p', 'o'},
			want:    String{v: "hippo", present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "String",
			give:    "camel",
			want:    String{v: "camel", present: true, initialized: true},
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
			got := &String{}
			if err := got.Scan(tt.give); (err != nil) != tt.wantErr {
				t.Errorf("String.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("String.Scan() = %v, want %v", got, tt.want)
			}
		})
	}
}
