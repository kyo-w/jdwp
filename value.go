// Copyright (C) 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jdwp

// Value is a generic value that can be one of the following types:
// • bool           • Char           • int            • int8
// • int16          • int32          • int64          • float32
// • float64        • ArrayID        • ClassLoaderID  • ClassObjectID
// • ObjectID       • StringID       • ThreadGroupID  • ThreadID
// • nil
type Value interface{}

// ValueSlice contains a set of values
type ValueSlice []Value

// untaggedValue can hold the same types as Value, but when encoded / decoded it
// is not prefixed with a type tag.
type untaggedValue interface{}
