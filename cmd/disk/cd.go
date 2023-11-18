/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package disk

import (
	"strings"

	"github.com/halalcloud/golang-sdk/utils"
	"github.com/spf13/cobra"
)

// CdCmd represents the cd command
var CdCmd = &cobra.Command{
	Use:   "cd",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		currentDir := utils.GetCurrentDir()
		if len(args) == 0 {
			cmd.Run(PwdCmd, []string{})
			return
		}
		if args[0] == ".." {
			if currentDir == "/" {
				cmd.Run(PwdCmd, []string{})
				return
			}
			strs := strings.Split(currentDir, "/")
			strs = strs[:len(strs)-1]
			currentDir = strings.Join(strs, "/")
			if len(currentDir) == 0 {
				currentDir = "/"
			}
			utils.SetCurrentDir(currentDir)
			cmd.Run(PwdCmd, []string{})
		}
		if args[0] == "." {
			cmd.Run(PwdCmd, []string{})
		}
		if args[0] == "/" {
			currentDir = "/"
			utils.SetCurrentDir(currentDir)
			cmd.Run(PwdCmd, []string{})
		}
		if args[0] != "." && args[0] != ".." && args[0] != "/" {
			currentDir = currentDir + "/" + args[0]
			if strings.HasPrefix(currentDir, "//") {
				currentDir = currentDir[1:]
			}
			utils.SetCurrentDir(currentDir)
			cmd.Run(PwdCmd, []string{})
		}
	},
}

func init() {
	DiskCmd.AddCommand(CdCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cdCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cdCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
