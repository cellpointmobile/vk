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

	"github.com/cellpointmobile/vk/programs"
	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall a tool.",
	Long:  `Uninstall the given tool.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		progname := args[0]
		progs := programs.LoadPrograms(cmd.Flag("bindir").Value.String())
		if prog, ok := progs[progname]; ok {
			if prog.IsInstalled() {
				f := prog.GetFullPath()
				os.Remove(f)
				fmt.Printf("%s has been uninstalled.\n", progname)
			} else {
				fmt.Printf("%s is not installed.\n", progname)
			}
		} else {
			fmt.Printf("Unknown program: %s\n", progname)
		}

	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
