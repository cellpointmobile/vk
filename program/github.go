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
	"context"
	"errors"
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
	TagName     string // Optional, will be used to find release. Ex: kustomize will find kustomize/v3.3.0. Used when multiple programs are released from the same repo.
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

func findAsset(r []*github.ReleaseAsset, name string) (*github.ReleaseAsset, error) {
	for _, x := range r {
		if *x.Name == name {
			return x, nil
		}
	}
	return nil, errors.New("can't find asset")
}

func getLatestTag(ctx context.Context, client *github.Client, githubOwner string, githubRepo string) string {
	tags, _, err := client.Repositories.ListTags(ctx, githubOwner, githubRepo, &github.ListOptions{})
	if _, ok := err.(*github.RateLimitError); ok {
		fmt.Println("Github rate limit hit, please add personal API token.")
		return ""
	}
	var filteredTags []*string
	for _, tag := range tags {
		if strings.HasPrefix(*tag.Name, "kustomize") {
			filteredTags = append(filteredTags, tag.Name)
		}
	}
	return *filteredTags[0]
}

// GetLatestVersion returns the latest version available
func (p *GithubProgram) GetLatestVersion() (string, string, error) {
	var err error
	var r *github.RepositoryRelease
	var u string
	var tag string
	var v string
	client, ctx := NewGithubClient()
	releases, _, err := client.Repositories.ListReleases(ctx, p.GithubOwner, p.GithubRepo, &github.ListOptions{})
	if _, ok := err.(*github.RateLimitError); ok {
		fmt.Println("Github rate limit hit, please add personal API token.")
		return "", "", err
	}
	if p.TagName != "" {
		tag = getLatestTag(ctx, client, p.GithubOwner, p.GithubRepo)
		fmt.Printf("Latest tag: %s\n", tag)
	}
	for _, release := range releases {
		if *release.Prerelease == p.PreRelease {
			if tag == "" {
				r = release
				break
			} else {
				if *release.TagName == tag {
					r = release
					break
				}
			}
		}
	}
	if tag == "" {
		v = strings.TrimPrefix(r.GetTagName(), "v")
	} else {
		v = strings.TrimPrefix(tag, p.TagName+"/")
	}
	rx := strings.NewReplacer("{VERSION}", v)
	if p.DownloadURL == "" {
		rn := rx.Replace(p.ReleaseName)
		la, _, err := client.Repositories.ListReleaseAssets(ctx, p.GithubOwner, p.GithubRepo, *r.ID, &github.ListOptions{})
		if _, ok := err.(*github.RateLimitError); ok {
			fmt.Println("Github rate limit hit, please add personal API token.")
			return "", "", err
		}
		a, err := findAsset(la, rn)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error finding asset.")
			os.Exit(200)
		}
		u = a.GetBrowserDownloadURL()
	} else {
		u = rx.Replace(p.DownloadURL)
	}
	return v, u, err
}

// DownloadLatestVersion downloads the latest release and puts it into the bindir
func (p *GithubDirectDownloadProgram) DownloadLatestVersion() string {
	f := filepath.Join(p.Path, p.Cmd)
	v, url, err := p.GetLatestVersion()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't get latest version.")
		os.Exit(10)
	}
	bak := f + ".bak"
	os.Rename(f, bak)
	_, err = grab.Get(f, url)
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
	v, url, err := p.GetLatestVersion()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't get latest version.")
		os.Exit(10)
	}
	rx := strings.NewReplacer("{VERSION}", v)
	err = file.ExtractFromTar(
		url,
		rx.Replace(p.Filename),
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
	v, url, err := p.GetLatestVersion()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't get latest version.")
		os.Exit(10)
	}
	rx := strings.NewReplacer("{VERSION}", v)
	err = file.ExtractFromZip(
		url,
		rx.Replace(p.Filename),
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
