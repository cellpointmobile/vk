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

package programs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/drzero42/vk/program"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/tidwall/gjson"
)

// LoadPrograms returns a map of programs
func LoadPrograms(bindir string) map[string]program.IProgram {
	path := os.ExpandEnv(bindir)
	cacheclient := httpcache.NewTransport(diskcache.New(os.ExpandEnv("$HOME/.vk/definitions-cache"))).Client()
	resp, err := cacheclient.Get("https://drzero42.github.io/vk-definitions/vk-definitions.json")
	if err != nil {
		panic(fmt.Sprintf("Could not download definitions: %s\n", err))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("Could not download definitions: %s\n", err))
	}

	directdownload := gjson.GetBytes(body, "github.directdownload")
	untarfile := gjson.GetBytes(body, "github.untarfile")
	unzipfile := gjson.GetBytes(body, "github.unzipfile")
	hashicorp := gjson.GetBytes(body, "hashicorp")

	progs := make(map[string]program.IProgram)

	for _, v := range directdownload.Array() {
		var prog program.GithubDirectDownloadProgram
		err := json.Unmarshal([]byte(v.String()), &prog)
		if err != nil {
			panic(err)
		}
		prog.Command.Path = path
		progs[prog.Command.Cmd] = &prog
	}
	for _, v := range untarfile.Array() {
		var prog program.GithubDownloadUntarFileProgram
		err := json.Unmarshal([]byte(v.String()), &prog)
		if err != nil {
			panic(err)
		}
		prog.Command.Path = path
		progs[prog.Command.Cmd] = &prog
	}
	for _, v := range unzipfile.Array() {
		var prog program.GithubDownloadUnzipFileProgram
		err := json.Unmarshal([]byte(v.String()), &prog)
		if err != nil {
			panic(err)
		}
		prog.Command.Path = path
		progs[prog.Command.Cmd] = &prog
	}
	for _, v := range hashicorp.Array() {
		var prog program.HashicorpProgram
		err := json.Unmarshal([]byte(v.String()), &prog)
		if err != nil {
			panic(err)
		}
		prog.Command.Path = path
		progs[prog.Command.Cmd] = &prog
	}

	return progs
}
