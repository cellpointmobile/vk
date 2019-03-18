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
	"net/http"
	"os"
	"sort"

	"github.com/cellpointmobile/vk/program"

	"github.com/cellpointmobile/vk/programs"
	"github.com/spf13/cobra"
)

func debugProgram(p program.IProgram) {
	fmt.Printf("Debugging tool %s\n", p.GetCmd())
	fmt.Printf("Struct: %#v\n", p)
	isInstalled := p.IsInstalled()
	fmt.Printf("Is installed: %t\n", isInstalled)
	if isInstalled {
		fmt.Printf("Local version: %s\n", p.GetLocalVersion())
	}
	v, url, err := p.GetLatestVersion()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't get latest version.")
		os.Exit(10)
	}
	fmt.Printf("Latest version: %s\n", v)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Something went wrong with the HTTP client: %s\n", err)
		os.Exit(20)
	}
	if resp.StatusCode == 200 {
		fmt.Printf("Download URL: %s\n", url)
	} else {
		fmt.Printf("Invalid DownloadURL: %s (Status code %d)\n", url, resp.StatusCode)
	}
}

// debugCmd represents the debug command
var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Debug a tool definition",
	Long:  `This subcommand debugs a tool definition.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		progs := programs.LoadPrograms(cmd.Flag("bindir").Value.String())
		if len(args) == 0 {
			keys := make([]string, 0, len(progs))
			for k := range progs {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				debugProgram(progs[k])
			}
		} else {
			progname := args[0]
			if prog, ok := progs[progname]; ok {
				debugProgram(prog)
			} else {
				fmt.Fprintf(os.Stderr, "Unknown program: %s\n", progname)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(debugCmd)
}
