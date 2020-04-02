package cmd

import (
	"fmt"
	"runtime"
	"github.com/spf13/cobra"

	"github.com/murer/lhproxy/server"

	. "github.com/murer/lhproxy/util"
)

var rootCmd *cobra.Command

func Config() {
	rootCmd = &cobra.Command{
		Use: "lhproxy", Short: "Last Hope Proxy",
		Version: fmt.Sprintf("%s-%s:%s", runtime.GOOS, runtime.GOARCH, Version),
	}

	rootCmd.AddCommand(&cobra.Command{
		Use: "version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf(rootCmd.Version)
			return nil
		},
	})

	configServer()
}

func configServer() {
	rootCmd.AddCommand(&cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) error {
			server.Start()
			return nil
		},
	})
}

func Execute() {
	err := rootCmd.Execute()
	Check(err)
}
