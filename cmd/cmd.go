package cmd

import (
	"fmt"
	"os"
	"runtime"
	"github.com/spf13/cobra"

	"github.com/murer/lhproxy/pipe"
	"github.com/murer/lhproxy/sockets"
	"github.com/murer/lhproxy/server"

	"github.com/murer/lhproxy/util"
)

var rootCmd *cobra.Command
var clientCmd *cobra.Command
var pipeCmd *cobra.Command

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

	clientCmd = &cobra.Command{Use:"client"}
	rootCmd.AddCommand(clientCmd)

	pipeCmd = &cobra.Command{Use:"pipe"}
	clientCmd.AddCommand(pipeCmd)

	configPipeNative()
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

func configPipeNative() {
	pipeCmd.AddCommand(&cobra.Command{
		Use: "native <host>:<port>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p := &pipe.Pipe{
				Scks: sockets.GetNative(),
				Address: args[0],
				Reader: os.Stdin,
				Writer: os.Stdout,
			}
			p.Execute()
			return nil
		},
	})
}

func Execute() {
	err := rootCmd.Execute()
	util.Check(err)
}
