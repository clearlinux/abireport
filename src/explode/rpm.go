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
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// RPM will explode all RPMs passed to it and return the path to
// the "root" to walk.
func RPM(pkgs []string) (string, error) {
	rootDir, err := ioutil.TempDir("", "abireport-rpm")
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
		}...)
		// Pipe rpm into cpio
		r, w := io.Pipe()
		defer r.Close()
		rpm.Stdout = w
		cpio.Stdin = r
		cpio.Stdout = nil
		cpio.Stderr = os.Stderr
		cpio.Dir = rootDir

		rpm.Start()
		cpio.Start()
		go func() {
			defer w.Close()
			rpm.Wait()
		}()
		if err := cpio.Wait(); err != nil {
			r.Close()
			return "", err
		}
		r.Close()
	}

	return rootDir, nil
}
