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
	"os/exec"
	"path/filepath"
)

// RPM will explode all RPMs passed to it and return the path to
// the "root" to walk.
func RPM(pkgs []string) (string, error) {
	rootDir, err := ioutil.TempDir("/var/tmp", "abireport-rpm")
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
		rpm := exec.Command("rpm2cpio", []string{
			fp,
		}...)
		cpio := exec.Command("cpio", []string{
			"-i",
			"-m",
			"-d",
			"--quiet",
			"-u",
		}...)

		cpio.Stdin, _ = rpm.StdoutPipe()
		cpio.Stdout = nil
		cpio.Stderr = nil
		cpio.Dir = rootDir
		cpio.Start()
		if err := rpm.Run(); err != nil {
			return "", err
		}
		if err := cpio.Wait(); err != nil {
			return "", err
		}
	}

	return rootDir, nil
}
