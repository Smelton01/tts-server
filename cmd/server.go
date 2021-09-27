/*
Copyright © 2021 Simon Mduduzi Juba scimail09@gmail.com
*/
package cmd

import (
	"github.com/smelton01/tts-server/api"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start the gRPC server",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		api.Serve()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
