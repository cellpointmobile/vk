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
	"os"
	"path/filepath"

	"github.com/drzero42/vk/file"
	"github.com/hashicorp/go-checkpoint"
)

// HashicorpProgram is for Hashicorp programs
type HashicorpProgram struct {
	Command
}

// GetLatestVersion returns the latest version number available
func (p *HashicorpProgram) GetLatestVersion() (string, error) {
	c, err := checkpoint.Check(&checkpoint.CheckParams{
		Product: "terraform",
	})
	return c.CurrentVersion, err
}

// GetLatestDownloadURL return the URL to download the latest version from
func (p *HashicorpProgram) GetLatestDownloadURL() string {
	c, err := checkpoint.Check(&checkpoint.CheckParams{
		Product: "terraform",
	})
	if err != nil {
		panic(err)
	}
	v := c.CurrentVersion
	url := c.CurrentDownloadURL + p.Cmd + "_" + v + "_linux_amd64.zip"
	return url
}

// DownloadLatestVersion downloads and extracts the latest version
func (p *HashicorpProgram) DownloadLatestVersion() string {
	f := filepath.Join(p.Path, p.Cmd)
	v, err := p.GetLatestVersion()
	if err != nil {
		panic("Can't get latest version.")
	}
	err = file.ExtractFromZip(
		p.GetLatestDownloadURL(),
		p.Cmd,
		f)
	if err != nil {
		panic(err)
	}
	if err = os.Chmod(f, 0755); err != nil {
		panic(err)
	}
	return v
}

// NewHashicorpProgram returns a new HashicorpProgram
func NewHashicorpProgram(
	cmd string,
	path string,
	versionRegexp string) *HashicorpProgram {
	prog := &HashicorpProgram{
		Command: Command{
			Cmd:           cmd,
			Path:          path,
			VersionArg:    "version",
			VersionRegexp: versionRegexp,
		},
	}
	return prog
}
