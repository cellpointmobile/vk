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

package programs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/cellpointmobile/vk/program"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

// LoadPrograms returns a map of programs
func LoadPrograms(bindir string) map[string]program.IProgram {
	path := os.ExpandEnv(bindir)
	url := viper.GetString("definitions")
	var d []byte
	var err error
	definitionsCache := os.ExpandEnv("$HOME/.vk/definitions-cache")
	if program.ClearCache {
		os.RemoveAll(definitionsCache)
	}
	if strings.HasPrefix(url, "http") {
		cacheclient := httpcache.NewTransport(diskcache.New(definitionsCache)).Client()
		resp, err := cacheclient.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not download definitions: %s\n", err)
			os.Exit(40)
		}
		defer resp.Body.Close()
		d, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not download definitions: %s\n", err)
			os.Exit(40)
		}
	} else {
		d, err = ioutil.ReadFile(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading definitions: %s\n", err)
			os.Exit(120)
		}
	}

	directdownload := gjson.GetBytes(d, "github.directdownload")
	untarfile := gjson.GetBytes(d, "github.untarfile")
	unzipfile := gjson.GetBytes(d, "github.unzipfile")
	hashicorp := gjson.GetBytes(d, "hashicorp")

	progs := make(map[string]program.IProgram)

	for _, v := range directdownload.Array() {
		var prog program.GithubDirectDownloadProgram
		err := json.Unmarshal([]byte(v.String()), &prog)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not unmarshal definitions: %s\n", err)
			os.Exit(50)
		}
		prog.Command.Path = path
		progs[prog.Command.Cmd] = &prog
	}
	for _, v := range untarfile.Array() {
		var prog program.GithubDownloadUntarFileProgram
		err := json.Unmarshal([]byte(v.String()), &prog)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not unmarshal definitions: %s\n", err)
			os.Exit(50)
		}
		prog.Command.Path = path
		progs[prog.Command.Cmd] = &prog
	}
	for _, v := range unzipfile.Array() {
		var prog program.GithubDownloadUnzipFileProgram
		err := json.Unmarshal([]byte(v.String()), &prog)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not unmarshal definitions: %s\n", err)
			os.Exit(50)
		}
		prog.Command.Path = path
		progs[prog.Command.Cmd] = &prog
	}
	for _, v := range hashicorp.Array() {
		var prog program.HashicorpProgram
		err := json.Unmarshal([]byte(v.String()), &prog)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not unmarshal definitions: %s\n", err)
			os.Exit(50)
		}
		prog.Command.Path = path
		progs[prog.Command.Cmd] = &prog
	}

	return progs
}
