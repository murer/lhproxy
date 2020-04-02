package cmd

import (
	"fmt"
	"os"
	"runtime"
	"github.com/spf13/cobra"

	"github.com/murer/lhproxy/client"
	"github.com/murer/lhproxy/server"

	"github.com/murer/lhproxy/util"
)

var rootCmd *cobra.Command

func Config() {
	rootCmd = &cobra.Command{
		Use: "lhproxy", Short: "Last Hope Proxy",
		Version: fmt.Sprintf("%s-%s:%s", runtime.GOOS, runtime.GOARCH, util.Version),
	}

	rootCmd.AddCommand(&cobra.Command{
		Use: "version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf(rootCmd.Version)
			return nil
		},
	})

	configServer()
	configClient()
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

func configClient() {
	clientCmd := &cobra.Command{Use:"client"}
	rootCmd.AddCommand(clientCmd)

	clientCmd.AddCommand(&cobra.Command{
		Use: "pipe <lhproxy:port> <host>:<port>",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			pipe := &client.Pipe{
				RAddress: args[1],
				LHAddress: args[0],
				LReader: os.Stdin,
				LWriter: os.Stdout,
			}
			pipe.Execute()
			return nil
		},
	})
}

func Execute() {
	err := rootCmd.Execute()
	util.Check(err)
}
