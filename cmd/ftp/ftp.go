/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package ftp

import (
	"github.com/spf13/cobra"
)

// shareCmd represents the share commandshareCmd
var FtpCmd = &cobra.Command{
	Use:     "ftp",
	Aliases: []string{"sftp"},
	Short:   "ftp functions",
	Long:    `ftp interface for the HalalCloud API.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		// viper.BindPFlag()
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {

	FtpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// shareCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// shareCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
