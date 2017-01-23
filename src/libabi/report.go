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
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

var (
	// KnownExtensions is a set of filename extensions applied by the
	// Architecture's GetPathSuffix function. We use this in our
	// TruncateAll function.
	KnownExtensions = []string{
		"",
		"32",
	}

	// ReportOutputDir is where report files will be dumped to. This
	// is set to the current working directory by default.
	ReportOutputDir = "."
)

// TruncateAll will truncate all files matching the current prefix
// with all known extensions in abireport. We do this to ensure that
// missing libs & reports are made obvious in git diffs.
func TruncateAll(prefix string) error {
	for _, ext := range KnownExtensions {
		p1 := filepath.Join(ReportOutputDir, fmt.Sprintf("%ssymbols%s", prefix, ext))
		p2 := filepath.Join(ReportOutputDir, fmt.Sprintf("%sused_libs%s", prefix, ext))

		if err := truncateFile(p1); err != nil {
			return err
		}
		if err := truncateFile(p2); err != nil {
			return err
		}
	}
	return nil
}

// PathExists is a simple test for path existence
func PathExists(p string) bool {
	if st, err := os.Stat(p); err == nil && st != nil {
		return true
	}
	return false
}

// truncateFile will truncate the file only if it already existed
func truncateFile(path string) error {
	if !PathExists(path) {
		return nil
	}
	fi, err := os.Create(path)
	if err != nil {
		return err
	}
	fi.Close()
	return nil
}

// writeSymbols will take care of writing out all the symbols provided by
// the given Architecture bucket, in a $soname:$symbol mapping, sorted first
// by soname, second by symbol.
func (a *Report) writeSymbols(prefix string, bucket *Architecture) error {
	suffix := bucket.GetPathSuffix()
	symbolsPath := filepath.Join(ReportOutputDir, fmt.Sprintf("%ssymbols%s", prefix, suffix))

	// Grab the sonames
	var sonames []string
	for key := range bucket.Symbols {
		sonames = append(sonames, key)
	}
	sort.Strings(sonames)

	if len(sonames) < 1 {
		// Truncate the file if it did exist.
		if err := truncateFile(symbolsPath); err != nil {
			return err
		}
		return nil
	}

	// The symbols file
	symsFi, err := os.Create(symbolsPath)
	if err != nil {
		return err
	}
	defer symsFi.Close()

	// Emit soname:symbol mapping
	for _, soname := range sonames {
		var foundSymbols []string
		for symbol := range bucket.Symbols[soname] {
			foundSymbols = append(foundSymbols, symbol)
		}
		sort.Strings(foundSymbols)
		for _, symbol := range foundSymbols {
			if _, err = fmt.Fprintf(symsFi, "%s:%s\n", soname, symbol); err != nil {
				return err
			}
		}
	}
	return nil
}

// writeDeps will write out a sorted list of soname's that this architecture
// bucket depends on.
// It will also filter out the names that are provided to generate a true
// reported based on DT_NEEDED requirements.
func (a *Report) writeDeps(prefix string, bucket *Architecture) error {
	suffix := bucket.GetPathSuffix()
	depsPath := filepath.Join(ReportOutputDir, fmt.Sprintf("%sused_libs%s", prefix, suffix))

	// Emit dependencies
	var depNames []string
	for nom := range bucket.Dependencies {
		// Skip provided
		if _, ok := bucket.Symbols[nom]; ok {
			continue
		}
		if _, ok := bucket.HiddenSymbols[nom]; ok {
			continue
		}
		depNames = append(depNames, nom)
	}
	sort.Strings(depNames)

	if len(depNames) < 1 {
		if err := truncateFile(depsPath); err != nil {
			return err
		}
		return nil
	}

	// The "used_libs" dependencies file
	depsFi, err := os.Create(depsPath)
	if err != nil {
		return err
	}
	defer depsFi.Close()

	for _, dep := range depNames {
		if _, err = fmt.Fprintf(depsFi, "%s\n", dep); err != nil {
			return err
		}
	}

	return nil
}

// Report will dump the report to the given writer for the specified
// machine configuration
func (a *Report) Report(prefix string, bucket *Architecture) error {
	if err := a.writeSymbols(prefix, bucket); err != nil {
		return err
	}

	return a.writeDeps(prefix, bucket)
}
