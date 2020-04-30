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
	"os/exec"
	"path/filepath"
	"strings"
)

// Eopkg will explode all eopkgs passed to it and return the path to
// the "root" to walk.
func Eopkg(pkgs []string) (string, error) {
	rootDir, err := ioutil.TempDir("/var/tmp", "abireport-eopkg")
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

		eopkg := exec.Command("unzip", []string{
			"-p",
			fp,
			"install.tar.xz",
		}...)
		tar := exec.Command("tar", []string{
			"-xJf",
			"-",
		}...)
		// Pipe eopkg into tar
		r, w := io.Pipe()
		defer r.Close()
		eopkg.Stdout = w
		tar.Stdin = r
		tar.Stdout = nil
		tar.Stderr = nil
		tar.Dir = rootDir

		eopkg.Start()
		tar.Start()
		go func() {
			defer w.Close()
			eopkg.Wait()
		}()
		if err := tar.Wait(); err != nil {
			r.Close()
			return "", err
		}
		r.Close()
	}

	return rootDir, nil
}
