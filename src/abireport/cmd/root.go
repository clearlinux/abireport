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
	"github.com/spf13/cobra"
	"libabi"
)

var (
	// Prefix is applied to the base of all report files to enable easier
	// integration.
	Prefix string
)

// RootCmd is the "default command" of abireport
var RootCmd = &cobra.Command{
	Use:   "abireport",
	Short: "Generate ABI reports for binary files",
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&Prefix, "prefix", "p", "", "Prefix for generated files")
	RootCmd.PersistentFlags().StringVarP(&libabi.ReportOutputDir, "output-dir", "D", ".", "Output directory for reports")
}
