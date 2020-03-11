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
	"github.com/clearlinux/abireport/libabi"
	"github.com/spf13/cobra"
	"os"
)

// RootCmd is the "default command" of abireport
var scanTreeCommand = &cobra.Command{
	Use:   "scan-tree [root]",
	Short: "Generate report from a filesystem tree",
	Long: `Examine the file tree beginning at [root] and generate an ABI report
for it. This is assumed to be the root directory with a compliant tree underneath
it containing /usr/lib, etc.`,
	Example: `
abireport scan-tree extractedRootfs/`,
	RunE: scanTree,
}

func init() {
	RootCmd.AddCommand(scanTreeCommand)
}

// scanTree is the CLI handler for "scan-tree".
func scanTree(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("scan-tree takes exactly one argument")
	}

	dir := args[0]

	if !libabi.PathExists(dir) {
		fmt.Fprintf(os.Stderr, "Directory %s doesn't exist\n", dir)
		os.Exit(1)
	}

	// Generate the ABI Report walker
	abi, err := libabi.NewReport(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initialising libabi: %v\n", err)
		os.Exit(1)
	}

	// Walk the exploded package tree
	if err = abi.Walk(); err != nil {
		fmt.Fprintf(os.Stderr, "Error walking rootfs: %v\n", err)
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

	return nil
}
