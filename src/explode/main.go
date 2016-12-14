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

// Package explode is a utility library for exploding various package
// types, so that abireport can traverse their root directory.
package explode

import (
	"strings"
)

var (
	// OutputDir is where the package implementation dumped
	// the packages to and extracted inside. This is automatically removed
	// at shutdown.
	OutputDir string
)

func init() {
	OutputDir = ""
}

// An Func is a function prototype for abireport extraction methods
type Func func(pkgs []string) (string, error)

var (
	// Impls is the valid set of packages understood by abireport
	Impls = map[string]Func{
		"*.rpm":   RPM,
		"*.eopkg": Eopkg,
		"*.deb":   Dpkg,
	}
)

// GetTypeForFilename will return the appropriate Impls key for the
// given input file, if it can be found.
func GetTypeForFilename(name string) string {
	if strings.HasSuffix(name, ".rpm") {
		return "*.rpm"
	}
	if strings.HasSuffix(name, ".eopkg") {
		return "*.eopkg"
	}
	if strings.HasSuffix(name, ".deb") {
		return "*.deb"
	}
	return ""
}

// ShouldSkipName is a utility to help with skipping any unwanted packages
func ShouldSkipName(name string) bool {
	if strings.HasSuffix(name, ".src.rpm") {
		return true
	}
	return false
}
