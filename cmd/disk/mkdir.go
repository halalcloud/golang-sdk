/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package disk

import (
	"context"
	"fmt"
	"os"
	"time"

	pubUserFile "github.com/city404/v6-public-rpc-proto/go/v6/userfile"
	"github.com/halalcloud/golang-sdk/auth"
	"github.com/halalcloud/golang-sdk/constants"
	"github.com/halalcloud/golang-sdk/pkg/print"
	"github.com/halalcloud/golang-sdk/utils"
	"github.com/spf13/cobra"
	"github.com/zzzhr1990/go-common-entity/userfile"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// mkdirCmd represents the mkdir command
var MkdirCmd = &cobra.Command{
	Use:     "mkdir",
	Short:   "A brief description of your command",
	Aliases: []string{"md"},
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		currentDir := utils.GetCurrentDir()
		serv, err := auth.NewAuthService(constants.AppID, constants.AppVersion, constants.AppSecret, "")
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(args) == 0 {
			fmt.Println("mkdir: missing operand")
			return
		}
		newDir := userfile.NewFormattedPath(utils.GetCurrentOpDir(args, 0)).GetPath()
		sp := print.Spinner(os.Stdout, "Create Directory [%s] ...", newDir)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		result, err := pubUserFile.NewPubUserFileClient(serv.GetGrpcConnection()).Create(ctx, &pubUserFile.File{
			// Parent: &pubUserFile.File{Path: currentDir},
			Path: newDir,
		})
		if err != nil {
			sp(false)
			status, ok := status.FromError(err)
			if ok {
				if status.Code() == codes.NotFound {
					fmt.Printf("Directory [%s] not found, back to root.\n", currentDir)
					utils.SetCurrentDir("/")
					return
				}
			}
			fmt.Println(err)
			return
		}
		sp(true)
		print.SuccessStatusEvent(os.Stdout, "Directory [%s] created.", result.Path)
	},
}

func init() {
	DiskCmd.AddCommand(MkdirCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mkdirCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mkdirCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
