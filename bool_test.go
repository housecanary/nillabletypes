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

func TestNewBool(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name string
		give bool
		want Bool
	}{
		{
			name: "True",
			give: true,
			want: Bool{v: true, present: true, initialized: true},
		},
		{
			name: "False",
			give: false,
			want: Bool{v: false, present: true, initialized: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBool(tt.give); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNilBool(t *testing.T) {
	tests := []struct {
		name string
		want Bool
	}{
		{"Nil", Bool{present: false, initialized: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NilBool(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NilBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBool_Bool(t *testing.T) {
	tests := []struct {
		name string
		give Bool
		want bool
	}{
		{
			name: "True",
			give: Bool{v: true, present: true, initialized: true},
			want: true,
		},
		{
			name: "False",
			give: Bool{v: false, present: true, initialized: true},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.give.Bool(); got != tt.want {
				t.Errorf("Bool.Bool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBool_Nil(t *testing.T) {
	tests := []struct {
		name string
		give Bool
		want bool
	}{
		{
			name: "Nil",
			give: Bool{present: false, initialized: true},
			want: true,
		},
		{
			name: "Not Nil",
			give: Bool{v: false, present: true, initialized: true},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.give.Nil(); got != tt.want {
				t.Errorf("Bool.Nil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBool_String(t *testing.T) {
	tests := []struct {
		name string
		give Bool
		want string
	}{
		{
			name: "True",
			give: Bool{v: true, present: true, initialized: true},
			want: "true",
		},
		{
			name: "False",
			give: Bool{v: false, present: true, initialized: true},
			want: "false",
		},
		{
			name: "Nil",
			give: Bool{v: false, present: false, initialized: true},
			want: "false",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.give.String(); got != tt.want {
				t.Errorf("Bool.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBool_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		give    []byte
		want    Bool
		wantErr bool
	}{
		{
			name:    "String",
			give:    toBytes("true"),
			want:    Bool{},
			wantErr: true,
		},
		{
			name:    "Floating Point Number",
			give:    toBytes(77.098),
			want:    Bool{},
			wantErr: true,
		},
		{
			name:    "Integer",
			give:    toBytes(int(1)),
			want:    Bool{},
			wantErr: true,
		},
		{
			name:    "Boolean (True)",
			give:    toBytes(true),
			want:    Bool{v: true, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Boolean (False)",
			give:    toBytes(false),
			want:    Bool{v: false, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Null",
			give:    toBytes(nil),
			want:    Bool{present: false, initialized: true},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &Bool{}
			if err := got.UnmarshalJSON(tt.give); (err != nil) != tt.wantErr {
				t.Errorf("Bool.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("Bool.UnmarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBool_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		give    Bool
		want    []byte
		wantErr bool
	}{
		{
			name:    "Null",
			give:    Bool{present: false, initialized: true},
			want:    toBytes(nil),
			wantErr: false,
		},
		{
			name:    "True",
			give:    Bool{v: true, present: true, initialized: true},
			want:    toBytes(true),
			wantErr: false,
		},
		{
			name:    "False",
			give:    Bool{v: false, present: true, initialized: true},
			want:    toBytes(false),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.give.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Bool.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bool.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBool_Value(t *testing.T) {
	tests := []struct {
		name    string
		give    Bool
		want    driver.Value
		wantErr bool
	}{
		{
			name:    "Nil",
			give:    Bool{present: false, initialized: true},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "True",
			give:    Bool{v: true, present: true, initialized: true},
			want:    true,
			wantErr: false,
		},
		{
			name:    "False",
			give:    Bool{v: false, present: true, initialized: true},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.give.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Bool.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bool.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBool_Scan(t *testing.T) {
	tests := []struct {
		name    string
		give    interface{}
		want    Bool
		wantErr bool
	}{
		{
			name:    "Nil",
			give:    nil,
			want:    Bool{present: false, initialized: true},
			wantErr: false,
		},
		{
			name:    "Bool (True)",
			give:    true,
			want:    Bool{v: true, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Bool (False)",
			give:    false,
			want:    Bool{v: false, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Int (Zero)",
			give:    int64(0),
			want:    Bool{v: false, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Int (Not Zero)",
			give:    int64(1),
			want:    Bool{v: true, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Float (Zero)",
			give:    float64(0.0),
			want:    Bool{v: false, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Float (Not Zero)",
			give:    float64(0.1),
			want:    Bool{v: true, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Byte Slice",
			give:    []byte{'0'},
			wantErr: true,
		},
		{
			name:    "String",
			give:    "true",
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
			got := &Bool{}
			if err := got.Scan(tt.give); (err != nil) != tt.wantErr {
				t.Errorf("Bool.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("Bool.Scan() = %v, want %v", got, tt.want)
			}
		})
	}
}
