/*
Copyright © 2023 Halal Cloud zzzhr@hotmail.com
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/halalcloud/golang-sdk/cmd/disk"
	"github.com/halalcloud/golang-sdk/cmd/user"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "halal",
	Short: "HalalCloud Golang SDK & CLI",
	Long:  `An Command Line Interface for HalalCloud API.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.halal-cloud.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(user.UserCmd)
	rootCmd.AddCommand(disk.DiskCmd)
	rootCmd.AddCommand(disk.ListCmd)
	rootCmd.AddCommand(disk.CdCmd)
	rootCmd.AddCommand(disk.PwdCmd)
	rootCmd.AddCommand(disk.MkdirCmd)
	rootCmd.AddCommand(disk.CreateCmd)
	rootCmd.AddCommand(disk.RmCmd)
	rootCmd.AddCommand(disk.CpCmd)
	rootCmd.AddCommand(disk.ListTrashCmd)
	rootCmd.AddCommand(disk.RnCmd)
	rootCmd.AddCommand(disk.SearchCmd)
	rootCmd.AddCommand(disk.DownloadCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".halal-cloud" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".halal-cloud")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "No config file found, Using default config.")
	}
}
