// Copyright © 2019 Anders Bruun Olsen <anders@bruun-olsen.net>
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

	"github.com/drzero42/vk/programs"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a tool",
	Long:  `Install latest version of the given tool.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		progname := args[0]
		progs := programs.LoadPrograms(cmd.Flag("bindir").Value.String())
		if prog, ok := progs[progname]; ok {
			if !prog.IsInstalled() || force {
				v := prog.DownloadLatestVersion()
				fmt.Printf("%s version %s has been installed.\n", progname, v)
			} else {
				fmt.Printf("%s is already installed.\n", progname)
			}
		} else {
			fmt.Printf("Unknown program: %s\n", progname)
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolVar(&force, "force", false, "Force installation of tool, overwriting installed version.")
}
