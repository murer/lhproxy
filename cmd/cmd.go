package cmd

import (
	"runtime"
	"github.com/spf13/cobra"
	. "github.com/murer/lhproxy/util"
)

var rootCmd *cobra.Command

func Config() {
	rootCmd = &cobra.Command{
		Use:     "lhproxy",
		Short:   "Last Hope Proxy",
		Version: "lhproxy " + runtime.GOOS + "-" + runtime.GOARCH + ":" + Version,
	}
}

func Execute() {
	err := rootCmd.Execute()
	Check(err)
}
