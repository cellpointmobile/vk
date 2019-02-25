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
	"context"
	"os"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/github"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

// IsLatestVersion returns true if installed program is latest version.
func IsLatestVersion(p IProgram) bool {
	if !p.IsInstalled() {
		return false
	}
	loV := p.GetLocalVersion()
	laV, err := p.GetLatestVersion()
	if err != nil {
		panic("Can't get latest version.")
	}
	localVersion := semver.MustParse(loV)
	latestVersion := semver.MustParse(laV)
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
	cacheclient := httpcache.NewTransport(diskcache.New(os.ExpandEnv("$HOME/.vk/github-cache"))).Client()
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
