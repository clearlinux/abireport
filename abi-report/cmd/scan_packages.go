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

package cmd

import (
	"fmt"
	"github.com/clearlinux/abireport/explode"
	"github.com/clearlinux/abireport/libabi"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

// RootCmd is the "default command" of abireport
var scanPkgsCommand = &cobra.Command{
	Use:   "scan-packages",
	Short: "Extract packages and generate report",
	Long: `Extract a set of packages, and create an ABI Report based
from their root directory.

If you do not pass a list of packages or a directory to scan, abireport
will look for them in the current directory.`,
	Run: scanPackages,
}

// Selected glob type
var exploderType string

func init() {
	RootCmd.AddCommand(scanPkgsCommand)
}

// scanPackages is the CLI handler for "scan-packages".
func scanPackages(cmd *cobra.Command, args []string) {

	// If no path is passed, use the current path
	var searchLocations []string
	if len(args) < 1 {
		searchLocations = append(searchLocations, ".")
	} else {
		searchLocations = append(searchLocations, args...)
	}

	extracts, err := locatePackages(searchLocations)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// i.e. passing a glob of *.src.rpm
	if len(extracts) < 1 {
		fmt.Fprintf(os.Stderr, "No usable packages found for extraction\n")
		os.Exit(1)
	}

	abi, err := explodeAndScan(extracts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in explode step: %v\n", err)
		os.Exit(1)
	}

	// Ensure we clean up existing reports
	if err = libabi.TruncateAll(Prefix); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot truncate existing reports: %v\n", err)
		os.Exit(1)
	}

	// Finally, create the report
	for _, arch := range abi.Arches {
		if err := abi.Report(Prefix, arch); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot generate report: %v\n", err)
			os.Exit(1)
		}
	}
}

// explodeAndScan will take care of exploding the packages, ensuring
// that a deferred cleanup always happen.
// Should it be successful, it will return a new ABI report that
// has walked the root.
func explodeAndScan(packages []string) (*libabi.Report, error) {
	exFunc := explode.Impls[exploderType]
	defer func() {
		if explode.OutputDir != "" && libabi.PathExists(explode.OutputDir) {
			os.RemoveAll(explode.OutputDir)
		}
	}()
	// Actually extract them
	root, err := exFunc(packages)
	if err != nil {
		return nil, err
	}

	// Generate the ABI Report walker
	abi, err := libabi.NewReport(root)
	if err != nil {
		return nil, err
	}

	// Walk the exploded package tree
	if err = abi.Walk(); err != nil {
		return nil, err
	}

	// Pass off to next step
	return abi, nil
}

// locatePackages will do the initial ground work of locating the
// packages.
func locatePackages(searchLocations []string) ([]string, error) {
	var discoveredPkgs []string

	for _, location := range searchLocations {
		found, err := discoverPackages(location)
		if err != nil {
			return nil, fmt.Errorf("error locating packages: %v", err)
		}
		discoveredPkgs = append(discoveredPkgs, found...)
	}

	var extracts []string
	for _, e := range discoveredPkgs {
		if explode.ShouldSkipName(e) {
			fmt.Fprintf(os.Stderr, "debug: Skipping '%s'\n", e)
			continue
		}
		extracts = append(extracts, e)
	}
	return extracts, nil
}

// discoverPackages will look at the given path and return either a new
// slice of found packages, or an error.
func discoverPackages(where string) ([]string, error) {
	st, err := os.Stat(where)
	if err != nil {
		return nil, err
	}

	// Search for a valid set of packages in the given directory
	if st.IsDir() {
		var paths []string
		var exp string
		for glo := range explode.Impls {
			fp := filepath.Join(where, glo)
			glob, _ := filepath.Glob(fp)
			if len(glob) > 0 {
				paths = append(paths, glob...)
				exp = glo
				break
			}
		}
		if len(paths) > 0 {
			// Ensure exploder type didn't change
			if exploderType != "" && exp != exploderType {
				return nil, fmt.Errorf("Cannot mix package type '%s' with '%s'", exp, exploderType)
			}
			exploderType = exp
			return paths, nil
		}
		return nil, fmt.Errorf("No packages in directory %s", where)
	}

	// Must be a file.
	exp := explode.GetTypeForFilename(where)
	if exp == "" {
		return nil, fmt.Errorf("No known package handler for %v", where)
	}
	// Don't allow mixing types..
	if exploderType != "" && exp != exploderType {
		return nil, fmt.Errorf("Cannot mix package type '%s' with '%s'", exp, exploderType)
	}
	exploderType = exp
	return []string{where}, nil
}
