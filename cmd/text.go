/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
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
		return tts.Read(strings.Join(args, " "), backend)
	},
}

func init() {
	rootCmd.AddCommand(textCmd)
}
