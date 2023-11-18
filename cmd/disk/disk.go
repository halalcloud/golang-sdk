/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package disk

import (
	"github.com/spf13/cobra"
)

// DiskCmd represents the disk command
var DiskCmd = &cobra.Command{
	Use:   "disk",
	Short: "Disk functions",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	// rootCmd.AddCommand(diskCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// diskCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	DiskCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
