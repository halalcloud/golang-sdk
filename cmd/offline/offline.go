/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package offline

import (
	"github.com/spf13/cobra"
)

// OfflineCmd represents the offline commandOfflineCmd
var OfflineCmd = &cobra.Command{
	Use:   "offline",
	Short: "offline functions",
	Long:  `offline interface for the HalalCloud API.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		// viper.BindPFlag()
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {

	OfflineCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// OfflineCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// OfflineCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
