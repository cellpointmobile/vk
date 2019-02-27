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

	"github.com/cavaliercoder/grab"
	"github.com/cellpointmobile/vk/file"
	"github.com/google/go-github/github"
)

// GithubProgram is a Program released via Github
type GithubProgram struct {
	GithubOwner string
	GithubRepo  string
	ReleaseName string // Will be appended when generating Github download URL. Ex: kustomize_{VERSION}_linux_amd64
	DownloadURL string // Optional, will be used instead of generating URL.
	PreRelease  bool   // Accept prereleases. Defaults to false
}

// GithubDirectDownloadProgram downloads a file directly
type GithubDirectDownloadProgram struct {
	Command
	GithubProgram
}

// GithubDownloadUntarFileProgram downloads a tarball and extracts a single file from it
type GithubDownloadUntarFileProgram struct {
	Command
	GithubProgram
	Filename string
}

// GithubDownloadUnzipFileProgram downloads a zip-file and extracts a single file from it
type GithubDownloadUnzipFileProgram struct {
	Command
	GithubProgram
	Filename string
}

// GetLatestVersion returns the latest version available
func (p *GithubProgram) GetLatestVersion() (string, error) {
	var v string
	var err error
	client, ctx := NewGithubClient()
	if p.PreRelease {
		listOptions := &github.ListOptions{PerPage: 1, Page: 1}
		releases, _, err := client.Repositories.ListReleases(ctx, p.GithubOwner, p.GithubRepo, listOptions)
		if _, ok := err.(*github.RateLimitError); ok {
			fmt.Println("Github rate limit hit, please add personal API token.")
			return "", err
		}
		v = releases[0].GetTagName()
	} else {
		latestRelease, _, err := client.Repositories.GetLatestRelease(ctx, p.GithubOwner, p.GithubRepo)
		if _, ok := err.(*github.RateLimitError); ok {
			fmt.Println("Github rate limit hit, please add personal API token.")
			return "", err
		}
		v = latestRelease.GetTagName()
	}
	return strings.TrimPrefix(v, "v"), err
}

// GetLatestDownloadURL returns the URL to download the latest release
func (p *GithubProgram) GetLatestDownloadURL() string {
	var url string
	version, err := p.GetLatestVersion()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't get latest version.")
		os.Exit(10)
	}
	r := strings.NewReplacer("{VERSION}", version,
		"{VVERSION}", "v"+version)
	if p.DownloadURL == "" {
		url = fmt.Sprintf("https://github.com/%s/%s/releases/download/{VVERSION}/%s", p.GithubOwner, p.GithubRepo, p.ReleaseName)
		url = r.Replace(url)
	} else {
		url = r.Replace(p.DownloadURL)
	}

	return url
}

// DownloadLatestVersion downloads the latest release and puts it into the bindir
func (p *GithubDirectDownloadProgram) DownloadLatestVersion() string {
	f := filepath.Join(p.Path, p.Cmd)
	v, err := p.GetLatestVersion()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't get latest version.")
		os.Exit(10)
	}
	bak := f + ".bak"
	os.Rename(f, bak)
	_, err = grab.Get(f, p.GetLatestDownloadURL())
	if err != nil {
		os.Rename(bak, f)
		fmt.Fprintf(os.Stderr, "Could not download update to %s: %s", p.GetCmd(), err)
		os.Exit(70)
	}
	if err = os.Chmod(f, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting chmod for downloaded file: %s", err)
		os.Exit(80)
	}
	os.Remove(bak)
	return v
}

// DownloadLatestVersion downloads and untars a file to the bindir
func (p *GithubDownloadUntarFileProgram) DownloadLatestVersion() string {
	f := filepath.Join(p.Path, p.Cmd)
	v, err := p.GetLatestVersion()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't get latest version.")
		os.Exit(10)
	}
	err = file.ExtractFromTar(
		p.GetLatestDownloadURL(),
		p.Filename,
		f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error extracting file from tarball: %s", err)
		os.Exit(90)
	}
	if err = os.Chmod(f, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting chmod for downloaded file: %s", err)
		os.Exit(80)
	}
	return v
}

// DownloadLatestVersion downloads and unzips a file to the bindir
func (p *GithubDownloadUnzipFileProgram) DownloadLatestVersion() string {
	f := filepath.Join(p.Path, p.Cmd)
	v, err := p.GetLatestVersion()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't get latest version.")
		os.Exit(10)
	}
	err = file.ExtractFromZip(
		p.GetLatestDownloadURL(),
		p.Filename,
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

// NewGithubDirectDownloadProgram returns a new GithubDirectDownloadProgram
func NewGithubDirectDownloadProgram(
	cmd string,
	path string,
	versionArg string,
	versionRegexp string,
	githubOwner string,
	githubRepo string,
	releaseName string,
	prerelease bool,
	downloadURL string) *GithubDirectDownloadProgram {
	prog := &GithubDirectDownloadProgram{
		Command: Command{
			Cmd:           cmd,
			Path:          path,
			VersionArg:    versionArg,
			VersionRegexp: versionRegexp,
		},
		GithubProgram: GithubProgram{
			GithubOwner: githubOwner,
			GithubRepo:  githubRepo,
			ReleaseName: releaseName,
			DownloadURL: downloadURL,
			PreRelease:  prerelease,
		},
	}
	return prog
}

// NewGithubDownloadUntarFileProgram returns a new GithubDirectDownloadProgram
func NewGithubDownloadUntarFileProgram(
	cmd string,
	path string,
	versionArg string,
	versionRegexp string,
	githubOwner string,
	githubRepo string,
	releaseName string,
	prerelease bool,
	downloadURL string,
	filename string) *GithubDownloadUntarFileProgram {
	prog := &GithubDownloadUntarFileProgram{
		Command: Command{
			Cmd:           cmd,
			Path:          path,
			VersionArg:    versionArg,
			VersionRegexp: versionRegexp,
		},
		GithubProgram: GithubProgram{
			GithubOwner: githubOwner,
			GithubRepo:  githubRepo,
			ReleaseName: releaseName,
			DownloadURL: downloadURL,
			PreRelease:  prerelease,
		},
		Filename: filename,
	}
	return prog
}

// NewGithubDownloadUnzipFileProgram returns a new GithubDirectDownloadProgram
func NewGithubDownloadUnzipFileProgram(
	cmd string,
	path string,
	versionArg string,
	versionRegexp string,
	githubOwner string,
	githubRepo string,
	releaseName string,
	prerelease bool,
	downloadURL string,
	filename string) *GithubDownloadUnzipFileProgram {
	prog := &GithubDownloadUnzipFileProgram{
		Command: Command{
			Cmd:           cmd,
			Path:          path,
			VersionArg:    versionArg,
			VersionRegexp: versionRegexp,
		},
		GithubProgram: GithubProgram{
			GithubOwner: githubOwner,
			GithubRepo:  githubRepo,
			ReleaseName: releaseName,
			DownloadURL: downloadURL,
			PreRelease:  prerelease,
		},
		Filename: filename,
	}
	return prog
}
