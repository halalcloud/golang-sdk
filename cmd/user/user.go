/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package user

import (
	"github.com/spf13/cobra"
)

// UserCmd represents the user command
var UserCmd = &cobra.Command{
	Use:   "user",
	Short: "user functions",
	Long:  `User interface for the HalalCloud API.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		// viper.BindPFlag()
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {

	UserCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
