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

func TestNewDate(t *testing.T) {
	tests := []struct {
		name string
		give string
		want Date
	}{
		{
			name: "Valid Format",
			give: "2019-11-12",
			want: Date{v: "2019-11-12", present: true, initialized: true},
		},
		{
			name: "Invalid Format",
			give: "2019/11/12",
			want: Date{v: "2019/11/12", present: true, initialized: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDate(tt.give); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNilDate(t *testing.T) {
	tests := []struct {
		name string
		want Date
	}{
		{"Nil", Date{v: "", present: false, initialized: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NilDate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NilDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_Nil(t *testing.T) {
	tests := []struct {
		name string
		give Date
		want bool
	}{
		{
			name: "Nil",
			give: Date{present: false, initialized: true},
			want: true,
		},
		{
			name: "Not Nil",
			give: Date{v: "", present: true, initialized: true},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.give.Nil(); got != tt.want {
				t.Errorf("Date.Nil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_DaysAgo(t *testing.T) {
	tests := []struct {
		name    string
		give    Date
		want    Int64
		wantErr bool
	}{
		{
			name:    "Nil",
			give:    Date{present: false, initialized: true},
			want:    Int64{},
			wantErr: false,
		},
		{
			name:    "Invalid Format",
			give:    Date{v: "2019/11/12", present: true, initialized: true},
			want:    Int64{},
			wantErr: true,
		},
		{
			name:    "One Day Ago",
			give:    Date{v: daysAgo(1), present: true, initialized: true},
			want:    Int64{v: 1, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "30 Days Ago",
			give:    Date{v: daysAgo(30), present: true, initialized: true},
			want:    Int64{v: 30, present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "Zero Days Ago",
			give:    Date{v: daysAgo(0), present: true, initialized: true},
			want:    Int64{v: 0, present: true, initialized: true},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.give.DaysAgo()
			if (err != nil) != tt.wantErr {
				t.Errorf("Date.DaysAgo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Date.DaysAgo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func daysAgo(days int) string {
	newTime := time.Now().AddDate(0 /* years */, 0 /* months */, -days)
	return newTime.Format("2006-01-02")
}

func TestDate_Before(t *testing.T) {
	tests := []struct {
		name  string
		this  Date
		other Date
		want  bool
	}{
		{
			name:  "this before other",
			this:  Date{v: "2023-05-20", present: true, initialized: true},
			other: Date{v: "2023-05-22", present: true, initialized: true},
			want:  true,
		},
		{
			name:  "this not before other",
			this:  Date{v: "2023-05-30", present: true, initialized: true},
			other: Date{v: "2023-05-22", present: true, initialized: true},
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.this.Before(tt.other)
			if err != nil {
				t.Errorf("Date.Before(): error = %v", err)
			}
			if got != tt.want {
				t.Errorf("Date.Before(): got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestDate_BeforeErr(t *testing.T) {
	tests := []struct {
		name  string
		this  Date
		other Date
	}{
		{
			name:  "nil",
			this:  Date{present: false, initialized: true},
			other: Date{v: "2023-05-22", present: true, initialized: true},
		},
		{
			name:  "malformed this",
			this:  Date{v: "20230525", present: true, initialized: true},
			other: Date{v: "2023-05-26", present: true, initialized: true},
		},
		{
			name:  "malformed other",
			this:  Date{v: "2023-05-25", present: true, initialized: true},
			other: Date{v: "20230526", present: true, initialized: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.this.Before(tt.other)
			if err == nil {
				t.Errorf("Date.Before(): got nil when expected error")
			}
		})
	}
}

func TestDate_After(t *testing.T) {
	tests := []struct {
		name  string
		this  Date
		other Date
		want  bool
	}{
		{
			name:  "this after other",
			this:  Date{v: "2023-05-30", present: true, initialized: true},
			other: Date{v: "2023-05-22", present: true, initialized: true},
			want:  true,
		},
		{
			name:  "this not after other",
			this:  Date{v: "2023-05-20", present: true, initialized: true},
			other: Date{v: "2023-05-22", present: true, initialized: true},
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.this.After(tt.other)
			if err != nil {
				t.Errorf("Date.After(): error = %v", err)
			}
			if got != tt.want {
				t.Errorf("Date.After(): got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestDate_AfterErr(t *testing.T) {
	tests := []struct {
		name  string
		this  Date
		other Date
	}{
		{
			name:  "nil",
			this:  Date{present: false, initialized: true},
			other: Date{v: "2023-05-22", present: true, initialized: true},
		},
		{
			name:  "malformed this",
			this:  Date{v: "20230525", present: true, initialized: true},
			other: Date{v: "2023-05-26", present: true, initialized: true},
		},
		{
			name:  "malformed other",
			this:  Date{v: "2023-05-25", present: true, initialized: true},
			other: Date{v: "20230526", present: true, initialized: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.this.After(tt.other)
			if err == nil {
				t.Errorf("Date.After(): got nil when expected error")
			}
		})
	}
}

func TestDate_String(t *testing.T) {
	tests := []struct {
		name string
		give Date
		want string
	}{
		{
			name: "Non-nil Empty",
			give: Date{v: "", present: true, initialized: true},
			want: "",
		},
		{
			name: "Nil",
			give: Date{v: "", present: false, initialized: true},
			want: "",
		},
		{
			name: "Valid Format",
			give: Date{v: "2019-11-12", present: true, initialized: true},
			want: "2019-11-12",
		},
		{
			name: "Invalid Format",
			give: Date{v: "11/12/2019", present: true, initialized: true},
			want: "11/12/2019",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.give.String(); got != tt.want {
				t.Errorf("Date.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		give    []byte
		want    Date
		wantErr bool
	}{
		{
			name:    "String (Valid Format)",
			give:    toBytes("2016-01-31"),
			want:    Date{v: "2016-01-31", present: true, initialized: true},
			wantErr: false,
		},
		{
			name:    "String (Invalid Format)",
			give:    toBytes("2016/01/31"),
			want:    Date{},
			wantErr: true,
		},
		{
			name:    "Floating Point Number",
			give:    toBytes(3.14159),
			want:    Date{},
			wantErr: true,
		},
		{
			name:    "Integer",
			give:    toBytes(100),
			want:    Date{},
			wantErr: true,
		},
		{
			name:    "Boolean",
			give:    toBytes(true),
			want:    Date{},
			wantErr: true,
		},
		{
			name:    "Null",
			give:    toBytes(nil),
			want:    Date{v: "", present: false, initialized: true},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &Date{}
			if err := got.UnmarshalJSON(tt.give); (err != nil) != tt.wantErr {
				t.Errorf("Date.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("Date.UnmarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		give    Date
		want    []byte
		wantErr bool
	}{
		{
			name:    "Present and Initialized (Valid Format)",
			give:    Date{v: "2019-11-12", present: true, initialized: true},
			want:    toBytes("2019-11-12"),
			wantErr: false,
		},
		{
			name:    "Present and Initialized (Invalid Format)",
			give:    Date{v: "2019/11/12", present: true, initialized: true},
			want:    toBytes("2019/11/12"),
			wantErr: false,
		},
		{
			name:    "Not Present",
			give:    Date{v: "2019-11-12", present: false, initialized: true},
			want:    toBytes(nil),
			wantErr: false,
		},
		{
			name:    "Not Initialized",
			give:    Date{v: "2019-11-12", present: true, initialized: false},
			want:    toBytes(nil),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.give.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Date.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Date.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_Value(t *testing.T) {
	tests := []struct {
		name    string
		give    Date
		want    driver.Value
		wantErr bool
	}{
		{
			name:    "Nil",
			give:    Date{present: false, initialized: true},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "Not Nil",
			give:    Date{v: "2018-07-11", present: true, initialized: true},
			want:    "2018-07-11",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.give.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Date.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Date.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_Scan(t *testing.T) {
	type args struct {
		src interface{}
	}
	tests := []struct {
		name    string
		give    interface{}
		want    Date
		wantErr bool
	}{
		{
			name:    "Nil",
			give:    nil,
			want:    Date{present: false, initialized: true},
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
			give:    toBytes("2019-10-01"),
			wantErr: true,
		},
		{
			name:    "String",
			give:    "2019-10-01",
			wantErr: true,
		},
		{
			name:    "Time",
			give:    time.Date(2011 /* year */, 12 /* month */, 12 /* day */, 0 /* hour */, 0 /* min */, 0 /* sec */, 0 /* nsec */, time.UTC),
			want:    Date{v: "2011-12-12", present: true, initialized: true},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &Date{}
			if err := got.Scan(tt.give); (err != nil) != tt.wantErr {
				t.Errorf("Date.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("Date.Scan() = %v, want %v", got, tt.want)
			}
		})
	}
}
