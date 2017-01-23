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

package explode

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Eopkg will explode all .eopkg's specified and then return
// the install/ path inside that exploded tree.
func Eopkg(pkgs []string) (string, error) {
	rootDir, err := ioutil.TempDir("", "abireport-eopkg")
	if err != nil {
		return "", err
	}
	// Ensure cleanup happens
	OutputDir = rootDir

	for _, archive := range pkgs {
		fp, err := filepath.Abs(archive)
		if err != nil {
			return "", err
		}
		// Don't want partials
		if strings.HasSuffix(archive, ".delta.eopkg") {
			continue
		}
		eopkg := exec.Command("uneopkg", []string{
			fp,
		}...)
		eopkg.Stdout = os.Stdout
		eopkg.Stderr = os.Stdout
		eopkg.Dir = rootDir

		fmt.Fprintf(os.Stderr, "Extracting %s\n", archive)
		if err = eopkg.Run(); err != nil {
			return "", err
		}
	}

	return filepath.Join(rootDir, "install"), nil
}
