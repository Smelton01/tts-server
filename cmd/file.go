/*
Copyright © 2021 Simon Mduduzi Juba scimail09@gmail.com
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/smelton01/tts-server/internal/tts"
	"github.com/spf13/cobra"
)

// fileCmd represents the file command
var fileCmd = &cobra.Command{
	Use:   "file [file-to-read-from]",
	Short: "read aloud from a file",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		b, err := ioutil.ReadFile(args[0])
		if err != nil {
			log.Fatalf("could not read file %s: %v", args[0], err)
		}
		message := string(b)
		audio, err := tts.Read(message, backend)
		if err != nil {
			return fmt.Errorf("could not read [%v]: %v", args[0], err)
		}
		return tts.PlayAudio(audio)
	},
}

func init() {
	rootCmd.AddCommand(fileCmd)
}
