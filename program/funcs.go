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
	"fmt"
	"os"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/github"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

// ClearCache variable for clear-cache flag
var ClearCache bool

// IsLatestVersion returns true if installed program is latest version.
func IsLatestVersion(p IProgram) bool {
	if !p.IsInstalled() {
		return false
	}
	loV := p.GetLocalVersion()
	laV, _, err := p.GetLatestVersion()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't get latest version.")
		os.Exit(10)
	}
	localVersion, err := semver.NewVersion(loV)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: Can not parse local version: '%s'\n", p.GetCmd(), loV)
		panic(err)
	}
	latestVersion, err := semver.NewVersion(laV)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: Can not parse latest version: '%s'\n", p.GetCmd(), laV)
		panic(err)
	}
	if localVersion.LessThan(latestVersion) {
		return false
	}
	return true
}

// NewGithubClient returns github.Client with auth if available otherwise unauthenticated
func NewGithubClient() (*github.Client, context.Context) {
	githubAPIToken := viper.GetString("github-api-token")
	var client *github.Client
	var ctx context.Context
	githubCache := os.ExpandEnv("$HOME/.vk/github-cache")
	if ClearCache {
		os.RemoveAll(githubCache)
	}
	cacheclient := httpcache.NewTransport(diskcache.New(githubCache)).Client()
	if githubAPIToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: githubAPIToken},
		)
		ctx = context.WithValue(context.Background(), oauth2.HTTPClient, cacheclient)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	} else {
		ctx = context.Background()
		client = github.NewClient(cacheclient)
	}
	return client, ctx
}
