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

// An Architecture is created for each ELF Machine type, and is used to group
// similar libraries and binaries in one place. This enables accounting for
// multilib/multiarch builds, to enable separate reports per architecture.
//
// This is required because we can enable different options per architecture
// in packaging, and base symbols will be different between them in almost all
// cases, i.e. ld-linux*
type Architecture struct {
	Machine       elf.Machine                // Corresponding machine for this configuration
	Symbols       map[string]map[string]bool // Symbols exported for this architecture
	HiddenSymbols map[string]map[string]bool // Symbols found but not exported
	Dependencies  map[string]bool            // Dependencies for this architecture
}

// NewArchitecture will create a new Architecture and initialise the fields
func NewArchitecture(m elf.Machine) *Architecture {
	return &Architecture{
		Machine:       m,
		Symbols:       make(map[string]map[string]bool),
		HiddenSymbols: make(map[string]map[string]bool),
		Dependencies:  make(map[string]bool),
	}
}

// GetSymbolsTarget will return the appropriate symbol store for the
// given record, based on it's symbol visibility (soname presence)
func (a *Architecture) GetSymbolsTarget(r *Record) map[string]map[string]bool {
	if r.Flags&RecordTypeExport == RecordTypeExport {
		return a.Symbols
	}
	return a.HiddenSymbols
}

// GetBucket will return an appropriate storage slot for the given
// record. If a bucket does not exist it will be created.
func (a *Report) GetBucket(record *Record) *Architecture {
	if arch, ok := a.Arches[record.Machine]; ok {
		return arch
	}

	bucket := NewArchitecture(record.Machine)
	a.Arches[record.Machine] = bucket
	return bucket
}

// GetPathSuffix will return an appropriate descriptor to use for the
// bucket configuration. This is used in the generated filenames
func (a *Architecture) GetPathSuffix() string {
	// TODO: Flesh out with other types.
	switch a.Machine {
	case elf.EM_X86_64:
		return ""
	case elf.EM_386:
		return "32"
	default:
		return a.Machine.String()
	}
}
