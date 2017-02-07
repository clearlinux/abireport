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
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// Dpkg will explode all .deb's specified and then return
// the install/ path inside that exploded tree.
func Dpkg(pkgs []string) (string, error) {
	rootDir, err := ioutil.TempDir("", "abireport-dpkg")
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
		dpkg := exec.Command("dpkg", []string{
			"-X",
			fp,
			filepath.Join(rootDir, "install"),
		}...)
		dpkg.Stdout = nil
		dpkg.Stderr = os.Stderr
		dpkg.Dir = rootDir

		if err = dpkg.Run(); err != nil {
			return "", err
		}
	}

	return filepath.Join(rootDir, "install"), nil
}
