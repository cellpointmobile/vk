vk installs and updates single-binary programs from Github and Hashicorp.

vk is an abbreviation of værktøjskasse *[ˈvæɐ̯gtʌjs-]*, which is a danish word
meaning toolbox.

Table of contents
=================
- [Table of contents](#table-of-contents)
- [Getting started](#getting-started)
- [Supported tools](#supported-tools)
- [Usage](#usage)
- [Github API rate limiting](#github-api-rate-limiting)
- [Problem abstract](#problem-abstract)
- [The solution](#the-solution)
- [Limitations](#limitations)
- [Tool definitions](#tool-definitions)

Getting started
===============
vk is made to install tools for a single user, not system-wide. It has, for
now, been hardcoded to use the directory `$HOME/.local/bin`. Before you start,
you have to ensure this directory is present and that you modify your PATH to
include it. When this has happened, go download the latest release of vk,
rename the binary to just vk, move it into the directory and give it execute
permission with `chmod +x vk`. Because vk is created just like the tools it
supports, it can keep itself updated (and even uninstall itself!).

Quick install:
```
curl -Lo ~/.local/bin/vk https://github.com/cellpointmobile/vk/releases/download/v0.2.2/vk_0.2.2_Linux_x86_64 && chmod +x ~/.local/bin/vk
```

Supported tools
===============
Here are some tools that are supported. The list is not exhaustive and can will
probably be expanded in the future.

* [Ark](https://github.com/heptio/velero/) (recently renamed to Velero)
* [Dive](https://github.com/wagoodman/dive/)
* [Docker-compose](https://github.com/docker/compose/)
* [Fly](https://github.com/concourse/concourse/)
* [Helm](https://github.com/helm/helm/)
* [Helmfile](https://github.com/roboll/helmfile/)
* [Img](https://github.com/genuinetools/img/)
* [K8sec](https://github.com/dtan4/k8sec)
* [Kops](https://github.com/kubernetes/kops)
* [Kustomize](https://github.com/kubernetes-sigs/kustomize/)
* [Minikube](https://github.com/kubernetes/minikube/)
* [Stern](https://github.com/wercker/stern/)
* [Terraform](https://www.terraform.io/)
* [terraform-docs](https://github.com/segmentio/terraform-docs/)
* [tflint](https://github.com/wata727/tflint/)

Usage
=====
To list tools that vk can install use the subcommand "available":
```
vk available
```

Optionally add the flag --all to also include already installed tools in the 
list.

To install a tool use the subcommand "install":
```
vk install minikube
```

To update all installed tools use the subcommand "update":
```
vk update
```

To install a specific tool only, also specify the tool in the "update" subcommand:
```
vk update minikube
```

To list installed tools use the subcommand "installed":
```
vk installed
```

To uninstall a tool use the subcommand "uninstall":
```
vk uninstall minikube
```

If you want to run update in a cronjob, there is a quiet flag you can use:
```
vk update --quiet
```

Update and install can also be forced to download and overwrite the local 
version with the latest version, even if the are the same:
```
vk [install|update] minikube --force
```

Bash/Zsh completion is available from the "completion" subcommand:
```
source <(vk completion [bash|zsh])
```

Github API rate limiting
========================
vk normally talks to the Github API unauthenticated, which means it is subject
to quite strict rate limiting - 60 requests per hour. If you run into this
problem (vk will tell you), you can either wait an hour or add a Github 
personal access token to your vk config. You can create a new token here:
https://github.com/settings/tokens/new. The token only needs the public_repo
scope and nothing else. When you have generated a token you have to insert the
following in `~/.vk/config.yaml` (create it, if it does not exist):
```
github-api-token: YOUR-TOKEN-HERE
```

Problem abstract
================
Infrastructure engineers and architects (the sysadmins of yesteryears) who are
working with cloud, containers, 12-factor etc need a bunch of modern tools to
do their jobs effectively. A lot of these tools have adopted a new paradigm
where they are released as single-binaries (often written in Go) on Github.
This makes for a fairly easy workflow when releasing, as there are nice tools
like goreleaser available that will build a binary, create a changelog and do
a release on Github for you.
For a single tool, this isn't too bad for a user - just go to the requisite
Github repo, download the latest release, optionally extract the binary from a
tarball or zipfile, copy the binary to a directory in your PATH and you are
ready to work. You then have to follow along on Github to be notified of new
releases and have to go through the same steps to update your tool. This
becomes quite tedious when you have more than three of these tools installed.
For the most part these tools aren't packaged for systems like apt, yum or
snap. And for the few tools that have unofficial packages, these are rarely
kept up-to-date. Quite a bit of these tools do have brew formulas available,
which does somewhat mitigate the problem. However not all tools have formulas
available, and brew for Linux requires a lot to run, while also mostly ending
up with your machine compiling the tools itself, if they even work.

The solution
============
vk is a simple tool that can download, extract, install, update and uninstall 
these types of tools. It does so by talking to the Github (and Hashicorps 
Checkpoint) API for information about the latest release, which it can check
against the version that is installed locally. For a user it is easy, 
convenient and fast to install and keep up-to-date these tools with vk.

Limitations
===========
vk is not a full blown package manager. It can not install specific versions
of tools, only the latest version. vk is also only for Linux AMD64 platforms
(for now).

Tool definitions
================
All tools that vk knows about are defined in a JSON file in the vk-definitions 
Github repo (https://github.com/cellpointmobile/vk-definitions).
