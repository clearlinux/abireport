//
// Copyright © Intel Corporation
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
	"fmt"
	"path/filepath"
	"strings"
)

// isLibraryDir will determine if a path is a valid library path worth
// exporting, i.e. /usr/lib64, etc.
func (a *Report) isLibraryDir(dir string) bool {
	for _, p := range a.libDirs {
		if dir == p {
			return true
		}
	}
	return false
}

// analyzeLibrary will examine the given shared library and populate
// the soname and symbols fields of the record
func (a *Report) analzyeLibrary(record *Record, file *elf.File) error {
	dynstring, err := file.DynString(elf.DT_SONAME)
	if err != nil {
		return err
	}

	// Export the symbol when we find a soname
	if len(dynstring) > 0 {
		record.Name = dynstring[0]
		record.Flags |= RecordTypeExport
	} else {
		record.Name = filepath.Base(record.Path)
	}

	// Unexport anything not in a library directory
	dirName := filepath.Dir(record.Path)
	if !a.isLibraryDir(dirName) && record.Flags&RecordTypeExport == RecordTypeExport {
		record.Flags ^= RecordTypeExport
	}

	symbols, err := file.DynamicSymbols()
	if err != nil {
		return err
	}

	nSections := elf.SectionIndex(len(file.Sections))

	// Try out best to emulate nm -g --defined-only --dynamic behaviour,
	// as used in autospec's older abireport.
	for _, sym := range symbols {
		// We only care for defined symbols
		// We want things in the .text section, and absolute symbols.
		inText := false
		sbind := elf.ST_BIND(sym.Info)
		// Skip weak symbols
		if sbind&elf.STB_WEAK == elf.STB_WEAK {
			continue
		}
		if sym.Section < nSections && file.Sections[sym.Section].Name == ".text" {
			inText = true
		}
		// If its not an absolute *and* its not in text, skip it too.
		if sym.Section&elf.SHN_ABS != elf.SHN_ABS && !inText {
			continue
		}
		symType := elf.ST_TYPE(sym.Info)
		// Skip STT_GNU_IFUNC
		if symType == elf.STT_LOOS {
			continue
		}
		// Skip unnamed ABI
		nom := strings.TrimSpace(sym.Name)
		if nom == "" {
			continue
		}
		record.Symbols = append(record.Symbols, nom)
	}
	return nil
}

// AnalyzeOne will attempt to analyze the given record, and store
// the appropriate details for a later report.
func (a *Report) AnalyzeOne(record *Record) error {
	file, err := elf.Open(record.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Determine if it's a shared library or executable
	if file.FileHeader.Type == elf.ET_DYN {
		record.Flags |= RecordTypeLibrary
	} else if file.FileHeader.Type == elf.ET_EXEC {
		record.Flags |= RecordTypeExecutable
	} else {
		return nil
	}

	// Analyze library elements
	if record.Flags&RecordTypeLibrary == RecordTypeLibrary {
		if err := a.analzyeLibrary(record, file); err != nil {
			return err
		}
	} else {
		// Set name to the executable name.
		record.Name = filepath.Base(record.Path)
	}

	// Set the appropriate class
	if file.FileHeader.Class == elf.ELFCLASS32 {
		record.Flags |= RecordType32bit
	} else if file.FileHeader.Class == elf.ELFCLASS64 {
		record.Flags |= RecordType64bit
	} else {
		return fmt.Errorf("Unknown ELF Class: %s", record.Path)
	}

	// Store the machine also
	record.Machine = file.FileHeader.Machine

	// Grab the required dependencies
	used, err := file.DynString(elf.DT_NEEDED)
	if err != nil {
		return err
	}

	record.Dependencies = used

	return nil
}
