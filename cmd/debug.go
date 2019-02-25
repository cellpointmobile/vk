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

// debugCmd represents the debug command
var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Debug a tool definition",
	Long:  `This subcommand debugs a tool definition.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		progname := args[0]
		progs := programs.LoadPrograms(cmd.Flag("bindir").Value.String())
		if prog, ok := progs[progname]; ok {
			fmt.Printf("Struct: %#v\n", prog)
			isInstalled := prog.IsInstalled()
			fmt.Printf("Is installed: %t\n", isInstalled)
			if isInstalled {
				fmt.Printf("Local version: %s\n", prog.GetLocalVersion())
			}
			v, err := prog.GetLatestVersion()
			if err != nil {
				panic("Can't get latest version.")
			}
			fmt.Printf("Latest version: %s\n", v)
			fmt.Printf("Download URL: %s\n", prog.GetLatestDownloadURL())
		} else {
			fmt.Printf("Unknown program: %s\n", progname)
		}
	},
}

func init() {
	rootCmd.AddCommand(debugCmd)
}
