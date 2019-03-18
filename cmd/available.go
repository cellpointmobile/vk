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
	"os"
	"sort"

	"github.com/cellpointmobile/vk/program"

	"github.com/cellpointmobile/vk/programs"
	"github.com/spf13/cobra"
)

// availableCmd represents the available command
var availableCmd = &cobra.Command{
	Use:   "available",
	Short: "List tools available for install",
	Long:  `Lists all available tools that are not already installed.`,
	Run: func(cmd *cobra.Command, args []string) {
		progs := programs.LoadPrograms(cmd.Flag("bindir").Value.String())
		all, _ := cmd.Flags().GetBool("all")
		keys := make([]string, 0, len(progs))
		for k := range progs {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			prog := progs[k]
			if all {
				v, _, err := prog.GetLatestVersion()
				if err != nil {
					fmt.Fprintln(os.Stderr, "Can't get latest version.")
					os.Exit(10)
				}
				fmt.Printf("%s version %s", prog.GetCmd(), v)
				if prog.IsInstalled() {
					if program.IsLatestVersion(prog) {
						fmt.Printf(" (installed)")
					} else {
						lv := prog.GetLocalVersion()
						fmt.Printf(" (%s installed)", lv)
					}
				}
				fmt.Printf("\n")
			} else {
				if !prog.IsInstalled() {
					v, _, err := prog.GetLatestVersion()
					if err != nil {
						fmt.Fprintln(os.Stderr, "Can't get latest version.")
						os.Exit(10)
					}
					fmt.Printf("%s version %s\n", prog.GetCmd(), v)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(availableCmd)
	availableCmd.Flags().Bool("all", false, "Include installed.")
}
