package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "Show the largest files in the given path",
	Long:  `Quickly Scan a directory anmd find large files`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("files called")
	},
}

func init() {
	rootCmd.AddCommand(filesCmd)
	
}
