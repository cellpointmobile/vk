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

var quiet bool

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update installed tools to latest version.",
	Long: `Go through all installed tools, look up the latest version available
	and update if the local version is not the latest.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		progs := programs.LoadPrograms(cmd.Flag("bindir").Value.String())
		if len(args) == 0 {
			keys := make([]string, 0, len(progs))
			for k := range progs {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				prog := progs[k]
				if prog.IsInstalled() {
					if !program.IsLatestVersion(prog) || force {
						v := prog.DownloadLatestVersion()
						if !quiet {
							fmt.Printf("Updating %s to version %s\n", prog.GetCmd(), v)
						}
					}
				}
			}
		} else {
			progname := args[0]
			if prog, ok := progs[progname]; ok {
				if prog.IsInstalled() {
					if !program.IsLatestVersion(prog) || force {
						v := prog.DownloadLatestVersion()
						if !quiet {
							fmt.Printf("Updating %s to version %s\n", prog.GetCmd(), v)
						}
					} else {
						if !quiet {
							fmt.Printf("%s is already latest version.\n", prog.GetCmd())
						}
					}
				} else {
					fmt.Fprintf(os.Stderr, "%s is not installed.", prog.GetCmd())
				}
			} else {
				fmt.Fprintf(os.Stderr, "Unknown program: %s\n", progname)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolVar(&force, "force", false, "Force installation of tool, overwriting installed version.")
	updateCmd.Flags().BoolVarP(&quiet, "quiet", "p", false, "Only output errors. Makes it suitable for cronjobs.")
}
