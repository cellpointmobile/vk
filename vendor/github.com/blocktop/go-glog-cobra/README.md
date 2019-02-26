# go-glog-cobra

[![standard-readme compliant](https://img.shields.io/badge/standard--readme-OK-green.svg?style=flat-square)](https://github.com/RichardLitt/standard-readme)

> Forms a bridge between [cobra](https://github.com/spf13/cobra) and [glog](https://godoc.org/github.com/golang/glog). Automatically defines cobra flags that match the glog flagset and passes values from cobra to [flag](https://godoc.org/flag) to satisfy the interface.

## Table of Contents

- [Install](#install)
- [Usage](#usage)
- [API](#api)
- [Maintainers](#maintainers)
- [Contribute](#contribute)
- [License](#license)

## Install

```
go get -u https://github.com/blocktop.go-glog-cobra
```

## Usage

```go
import glogcobra "github.com/blocktop.go-glog-cobra"

func init() {
	cobra.OnInitialize(initConfig)

	// rootCmd is defined by standard cobra configuration
	glogcobra.Init(rootCmd)
}

// initConfig is part of the standard rootCmd configuration that
// cobra creates.
func initConfig() {
	// This will also call flag.Parse() if you have not already.
	glogcobra.Parse(rootCmd) 
	
	// glog will now have all the flags it needs
}
```

## Maintainers

[@strobus](https://github.com/strobus)

## Contribute

PRs accepted.

Small note: If editing the README, please conform to the [standard-readme](https://github.com/RichardLitt/standard-readme) specification.

## License

MIT © 2018 J. Strobus White
