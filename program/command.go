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

package program

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// Command defines command, version args and regexp to find version number.
type Command struct {
	Path          string
	Cmd           string
	VersionArg    string
	VersionRegexp string
}

// GetLocalVersion runs Cmd with versionArg and finds version using versionRegexp, which it returns.
func (p *Command) GetLocalVersion() string {
	args := strings.Split(p.VersionArg, " ")
	cmd := filepath.Join(p.Path, p.Cmd)
	version := exec.Command(cmd, args...)
	versionOut, err := version.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing version: %s", err)
		os.Exit(60)
	}
	vr := regexp.MustCompile(p.VersionRegexp)
	match := vr.FindStringSubmatch(
		string(versionOut))
	return match[1]
}

// IsInstalled checks if command is installed and returns boolean
func (p *Command) IsInstalled() bool {
	if _, err := os.Stat(filepath.Join(p.Path, p.Cmd)); err == nil {
		return true
	}
	return false
}

// GetCmd returns the command name
func (p *Command) GetCmd() string {
	return p.Cmd
}

// GetFullPath returns the full path to the command
func (p *Command) GetFullPath() string {
	return filepath.Join(p.Path, p.Cmd)
}
