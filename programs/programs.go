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
	"os"

	"github.com/drzero42/vk/program"
)

// LoadPrograms returns a map of programs
func LoadPrograms(bindir string) map[string]program.IProgram {
	path := os.ExpandEnv(bindir)
	Progs := make(map[string]program.IProgram)

	Progs["ark"] = program.NewGithubDownloadUntarFileProgram(
		"ark",
		path,
		"version",
		"Version: (v.+?)\n",
		"heptio",
		"velero",
		"ark-{VVERSION}-linux-amd64.tar.gz",
		false,
		"",
		"ark",
	)

	Progs["dive"] = program.NewGithubDownloadUntarFileProgram(
		"dive",
		path,
		"version",
		"dive (.+?)\n",
		"wagoodman",
		"dive",
		"dive_{VERSION}_linux_amd64.tar.gz",
		true,
		"",
		"dive",
	)

	Progs["docker-compose"] = program.NewGithubDirectDownloadProgram(
		"docker-compose",
		path,
		"version",
		"docker-compose version (.+?),",
		"docker",
		"compose",
		"docker-compose-Linux-x86_64",
		false,
		"",
	)

	Progs["fly"] = program.NewGithubDirectDownloadProgram(
		"fly",
		path,
		"--version",
		"(.+?)\n",
		"concourse",
		"concourse",
		"fly_linux_amd64",
		false,
		"",
	)

	Progs["helm"] = program.NewGithubDownloadUntarFileProgram(
		"helm",
		path,
		"version --client",
		"SemVer:\"(.+?)\"",
		"kubernetes",
		"helm",
		"",
		false,
		"https://storage.googleapis.com/kubernetes-helm/helm-{VVERSION}-linux-amd64.tar.gz",
		"linux-amd64/helm",
	)

	Progs["helmfile"] = program.NewGithubDirectDownloadProgram(
		"helmfile",
		path,
		"--version",
		"version (v.+?)\n",
		"roboll",
		"helmfile",
		"helmfile_linux_amd64",
		false,
		"",
	)

	Progs["img"] = program.NewGithubDirectDownloadProgram(
		"img",
		path,
		"version",
		"version.+?(v.+?)\n",
		"genuinetools",
		"img",
		"img-linux-amd64",
		false,
		"",
	)

	Progs["k8sec"] = program.NewGithubDownloadUntarFileProgram(
		"k8sec",
		path,
		"version",
		"version (v.+?),",
		"dtan4",
		"k8sec",
		"k8sec-{VVERSION}-linux-amd64.tar.gz",
		false,
		"",
		"linux-amd64/k8sec",
	)

	Progs["kops"] = program.NewGithubDirectDownloadProgram(
		"kops",
		path,
		"version",
		"Version (.+?) ",
		"kubernetes",
		"kops",
		"kops-linux-amd64",
		false,
		"",
	)

	Progs["kustomize"] = program.NewGithubDirectDownloadProgram(
		"kustomize",
		path,
		"version",
		"KustomizeVersion:(.+?) ",
		"kubernetes-sigs",
		"kustomize",
		"kustomize_{VVERSION}_linux_amd64",
		false,
		"",
	)

	Progs["minikube"] = program.NewGithubDirectDownloadProgram(
		"minikube",
		path,
		"version",
		"version: (v.+?)\n",
		"kubernetes",
		"minikube",
		"minikube-linux-amd64",
		false,
		"",
	)

	Progs["stern"] = program.NewGithubDirectDownloadProgram(
		"stern",
		path,
		"--version",
		"version (.+?)\n",
		"wercker",
		"stern",
		"stern_linux_amd64",
		false,
		"",
	)

	Progs["terraform"] = program.NewHashicorpProgram(
		"terraform",
		path,
		"Terraform (v.+?)\n",
	)

	Progs["terraform-docs"] = program.NewGithubDirectDownloadProgram(
		"terraform-docs",
		path,
		"--version",
		"(.+?)\n",
		"segmentio",
		"terraform-docs",
		"terraform-docs-{VVERSION}-linux-amd64",
		false,
		"",
	)

	Progs["tflint"] = program.NewGithubDownloadUnzipFileProgram(
		"tflint",
		path,
		"--version",
		"version (.+?)\n",
		"wata727",
		"tflint",
		"tflint_linux_amd64.zip",
		false,
		"",
		"tflint",
	)

	return Progs
}
