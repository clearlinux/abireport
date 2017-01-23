//
// Copyright © 2016-2017 Intel Corporation
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
//

package libabi

import (
	"debug/elf"
)

// A RecordType is a flag which can be OR'd to define a type of Record
// encountered.
type RecordType uint8

const (
	// RecordTypeExecutable indicates the file was an executable
	RecordTypeExecutable RecordType = 1 << iota

	// RecordTypeLibrary indicates a shared library
	RecordTypeLibrary RecordType = 1 << iota

	// RecordType64bit indicates a 64-bit file
	RecordType64bit RecordType = 1 << iota

	// RecordType32bit indicates a 32-bit file
	RecordType32bit RecordType = 1 << iota

	// RecordTypeExport is set when we encounter valid ABI, i.e. a library
	// with a soname that we can export.
	RecordTypeExport RecordType = 1 << iota
)

// A Record is literally a recording of an encounter, with a file that
// we believe to hold some interest.
type Record struct {
	Path         string      // Where we found the file
	Flags        RecordType  // A bitwise set of RecordType
	Name         string      // Either the soname or the basename
	Dependencies []string    // DT_NEEDED dependencies
	Symbols      []string    // Dynamic defined symbols
	Machine      elf.Machine // Corresponding machine
}
