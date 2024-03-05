package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "Show the largest files in the given path",
	Long:  `Quickly Scan a directory anmd find large files`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("files called")
		if Debug {
			for key, value := range viper.GetViper().AllSettings() {
				log.WithFields(log.Fields{
					key: value,
				}).Info("Command Flag")
			}

		}
	},
}

func init() {
	rootCmd.AddCommand(filesCmd)
}
