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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PwdCmd represents the pwd command
var PwdCmd = &cobra.Command{
	Use:   "pwd",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		currentDir := utils.GetCurrentDir()
		if len(currentDir) == 0 {
			currentDir = "/"
		}
		serv, err := auth.NewAuthService(constants.AppID, constants.AppVersion, constants.AppSecret, "")
		if err != nil {
			fmt.Println(err)
			return
		}
		sp := print.Spinner(os.Stdout, "Get Directory Info [%s] ...", currentDir)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		result, err := pubUserFile.NewPubUserFileClient(serv.GetGrpcConnection()).Get(ctx, &pubUserFile.File{
			Path: currentDir,
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
		utils.SetCurrentDir(result.Path)
		fmt.Printf("Current Directory: %s\n\n", result.Path)
	},
}

func init() {
	DiskCmd.AddCommand(PwdCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pwdCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pwdCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
