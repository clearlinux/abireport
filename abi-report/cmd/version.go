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
	"github.com/spf13/cobra"
)

// ABIReportVersion is the public version of the tool. The value is initialized
// with a linker option to `go build`. See the toplevel Makefile for details.
var ABIReportVersion = "UNKNOWN"

// versionCmd handles "abireport version"
var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Print the version and exit",
	Long:  "Print the version & license information for abireport.",
	Run:   printVersion,
}

func init() {
	RootCmd.AddCommand(versionCommand)
}

// Print the application version and exit.
func printVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("abireport version %v\n\nCopyright © 2016-2017 Intel Corporation\n", ABIReportVersion)
	fmt.Printf("Licensed under the Apache License, Version 2.0\n")
}
