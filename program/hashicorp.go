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

package program

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cellpointmobile/vk/file"
	"github.com/hashicorp/go-checkpoint"
)

// HashicorpProgram is for Hashicorp programs
type HashicorpProgram struct {
	Command
}

// GetLatestVersion returns the latest version number available
func (p *HashicorpProgram) GetLatestVersion() (string, error) {
	cmd := p.GetCmd()
	cache := os.ExpandEnv("$HOME/.vk/checkpoint-cache/" + cmd)
	if ClearCache {
		os.RemoveAll(cache)
	}
	c, err := checkpoint.Check(&checkpoint.CheckParams{
		Product:   cmd,
		CacheFile: cache,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting version from Checkpoint: %s", err)
		os.Exit(100)
	}
	return c.CurrentVersion, err
}

// GetLatestDownloadURL return the URL to download the latest version from
func (p *HashicorpProgram) GetLatestDownloadURL() string {
	cmd := p.GetCmd()
	v, _ := p.GetLatestVersion()
	r := strings.NewReplacer(
		"{VERSION}", v,
		"{VVERSION}", "v"+v,
		"{CMD}", cmd)
	url := "https://releases.hashicorp.com/{CMD}/{VERSION}/{CMD}_{VERSION}_linux_amd64.zip"
	return r.Replace(url)
}

// DownloadLatestVersion downloads and extracts the latest version
func (p *HashicorpProgram) DownloadLatestVersion() string {
	f := filepath.Join(p.Path, p.Cmd)
	v, err := p.GetLatestVersion()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't get latest version.")
		os.Exit(10)
	}
	err = file.ExtractFromZip(
		p.GetLatestDownloadURL(),
		p.Cmd,
		f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error extracting file from zip: %s", err)
		os.Exit(90)
	}
	if err = os.Chmod(f, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting chmod for downloaded file: %s", err)
		os.Exit(80)
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
