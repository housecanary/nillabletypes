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

// +build !amd64,!arm64,!arm64be,!ppc64,!ppc64le,!mips64,!mips64le,!s390x,!sparc64

package nillabletypes

// Int represents a nil-able int
type Int = Int32

// NewInt makes a new non-nil Int
func NewInt(v int) Int {
	return NewInt32(int32(v))
}

// NilInt makes a new nil Int
func NilInt() Int {
	return NilInt32()
}

// Int returns the built-in int32 value
func (v Int) Int() int {
	return int(v.Int32())
}
