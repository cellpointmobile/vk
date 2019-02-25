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

// installedCmd represents the installed command
var installedCmd = &cobra.Command{
	Use:   "installed",
	Short: "List all installed tools",
	Long:  `Output a list of installed tools with their versions`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("The following programs are installed:")
		progs := programs.LoadPrograms(cmd.Flag("bindir").Value.String())
		for _, prog := range progs {
			if prog.IsInstalled() {
				fmt.Printf("%s: %s\n", prog.GetCmd(), prog.GetLocalVersion())
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(installedCmd)
}
