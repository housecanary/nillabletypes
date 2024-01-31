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
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/segmentio/encoding/json"
)

func TestNewFloat(t *testing.T) {
	type args struct {
		v float64
	}
	tests := []struct {
		name string
		give float64
		want Float
	}{
		{
			name: "Zero",
			give: 0,
			want: Float{v: 0, present: true, initialized: true},
		},
		{
			name: "Positive",
			give: 17.12098123,
			want: Float{v: 17.12098123, present: true, initialized: true},
		},
		{
			name: "Negative",
			give: -8.200001,
			want: Float{v: -8.200001, present: true, initialized: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFloat(tt.give); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNilFloat(t *testing.T) {
	tests := []struct {
		name string
		want Float
	}{
		{"Nil", Float{v: 0, present: false, initialized: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NilFloat(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NilFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat_Float(t *testing.T) {
	tests := []struct {
		name string
		give Float
		want float64
	}{
		{
			name: "Non-nil Zero",
			give: Float{v: 0, present: true, initialized: true},
			want: 0,
		},
		{
			name: "Nil",
			give: Float{v: 0, present: false, initialized: true},
			want: 0,
		},
		{
			name: "Positive",
			give: Float{v: 17.12098123, present: true, initialized: true},
			want: 17.12098123,
		},
		{
			name: "Negative",
			give: Float{v: -8.200001, present: true, initialized: true},
			want: -8.200001,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.give.Float(); got != tt.want {
				t.Errorf("Float.Float() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat_Nil(t *testing.T) {
	tests := []struct {
		name string
		give Float
		want bool
	}{
		{
			name: "Nil",
			give: Float{present: false, initialized: true},
			want: true,
		},
		{
			name: "Not Nil",
			give: Float{v: 0, present: true, initialized: true},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.give.Nil(); got != tt.want {
				t.Errorf("Float.Nil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat_String(t *testing.T) {
	tests := []struct {
		name string
		give Float
		want string
	}{
		{
			name: "Non-nil Zero",
			give: Float{v: 0, present: true, initialized: true},
			want: "0",
		},
		{
			name: "Nil",
			give: Float{v: 0, present: false, initialized: true},
			want: "0",
		},
		{
			name: "Positive",
			give: Float{v: 98.3091002, present: true, initialized: true},
			want: "98.3091002",
		},
		{
			name: "Negative",
			give: Float{v: -1000.0001, present: true, initialized: true},
			want: "-1000.0001",
		},
		{
			name: "Many Decimals",
			give: Float{v: 123.4560000000001, present: true, initialized: true},
			want: "123.4560000000001",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.give.String(); got != tt.want {
				t.Errorf("Float.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func toBytes(val interface{}) []byte {
	bytes, err := json.Marshal(val)
	if err != nil {
		fmt.Printf("Failed to convert value (%v) of type %T to bytes: %v\n", val, val, err)
	}
	return bytes
}

func TestFloat_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		give    []byte
		want    Float
		wantErr bool
	}{
		{
			name:    "String",
			give:    toBytes("3.14159"),
			wantErr: true,
		},
		{
			name:    "Floating Point Number",
			give:    toBytes(3.14159),
			want:    Float{v: 3.14159, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Integer",
			give:    toBytes(100),
			want:    Float{v: 100, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Boolean",
			give:    toBytes(true),
			want:    Float{},
			wantErr: true,
		},
		{
			name:    "Null",
			give:    toBytes(nil),
			want:    Float{v: 0, present: false, initialized: true},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &Float{}
			if err := got.UnmarshalJSON(tt.give); (err != nil) != tt.wantErr {
				t.Errorf("Float.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("Float.UnmarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		give    Float
		want    []byte
		wantErr bool
	}{
		{
			name:    "Present and Initialized",
			give:    Float{v: 7, present: true, initialized: true},
			want:    toBytes(7),
			wantErr: false,
		},
		{
			name:    "Not Present",
			give:    Float{v: 7, present: false, initialized: true},
			want:    toBytes(nil),
			wantErr: false,
		},
		{
			name:    "Not Initialized",
			give:    Float{v: 7, present: true, initialized: false},
			want:    toBytes(nil),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.give.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Float.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Float.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat_Value(t *testing.T) {
	tests := []struct {
		name    string
		give    Float
		want    driver.Value
		wantErr bool
	}{
		{
			name:    "Nil",
			give:    Float{present: false, initialized: true},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "Not Nil",
			give:    Float{v: 65.12, present: true, initialized: true},
			want:    65.12,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.give.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Float.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Float.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat_Scan(t *testing.T) {
	type args struct {
		src interface{}
	}
	tests := []struct {
		name    string
		give    interface{}
		want    Float
		wantErr bool
	}{
		{
			name:    "Nil",
			give:    nil,
			want:    Float{present: false, initialized: true},
			wantErr: false,
		},
		{
			name:    "Bool (True)",
			give:    true,
			want:    Float{v: 1, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Bool (False)",
			give:    false,
			want:    Float{v: 0, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Int",
			give:    int64(123),
			want:    Float{v: 123, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Float",
			give:    float64(123.456),
			want:    Float{v: 123.456, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Byte Slice (Valid)",
			give:    []byte{'3', '.', '1', '4'},
			want:    Float{v: 3.14, present: true, initialized: true},
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
			want:    Float{v: 92.96, present: true, initialized: true},
			wantErr: false,
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
			got := &Float{}
			if err := got.Scan(tt.give); (err != nil) != tt.wantErr {
				t.Errorf("Float.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("Float.Scan() = %v, want %v", got, tt.want)
			}
		})
	}
}
