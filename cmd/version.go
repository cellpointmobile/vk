// Copyright © 2019 Cellpoint Mobile
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	goVersion "github.com/christopherhein/go-version"

	"github.com/spf13/cobra"
)

var (
	json    = false
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Output version",
	Long:  `This subcommand outputs the version of vk.`,
	Run: func(_ *cobra.Command, _ []string) {
		var response string
		versionOutput := goVersion.New(version, commit, date)

		if json {
			response = versionOutput.ToJSON()
		} else {
			response = versionOutput.ToShortened()
		}
		fmt.Printf("%+v", response)
		return
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&json, "json", "j", false, "Output in JSON format.")
}
