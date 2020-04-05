package main

import (
	"github.com/murer/lhproxy/cmd"
	"github.com/murer/lhproxy/util"
)

var Version = "dev"

func main() {
	util.Version = Version
	cmd.Config()
	cmd.Execute()
}
