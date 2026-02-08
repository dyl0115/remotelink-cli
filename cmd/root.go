package cmd

import (
	"fmt"
	"os"
	"remotelink/config"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "remotelink",
	Short: "",
	Long:  "",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.LoadServers()
	},
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
