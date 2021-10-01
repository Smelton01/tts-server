/*
Copyright © 2021 Simon Mduduzi Juba scimail09@gmail.com
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/smelton01/tts-server/internal/tts"
	"github.com/spf13/cobra"
)

// textCmd represents the text command
var textCmd = &cobra.Command{
	Use:   "text [text-to-read]",
	Short: "read aloud a block of text.",
	Long:  `Long definition about text goes here.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		audio, err := tts.Read(strings.Join(args, " "), backend)
		if err != nil {
			return fmt.Errorf("could not read [%v]: %v", args[0], err)
		}
		return tts.PlayAudio(audio)
	},
}

func init() {
	rootCmd.AddCommand(textCmd)
}
