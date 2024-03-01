/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package share

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	pubFileShare "github.com/city404/v6-public-rpc-proto/go/v6/fileshare"
	"github.com/halalcloud/golang-sdk/auth"
	"github.com/halalcloud/golang-sdk/constants"
	"github.com/halalcloud/golang-sdk/pkg/print"
	"github.com/halalcloud/golang-sdk/utils"
	"github.com/spf13/cobra"
	"github.com/zzzhr1990/go-common-entity/userfile"
)

// addCmd represents the info command
var addCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"add"},
	Short:   "create share",
	Long:    `create a share`,
	Run: func(cmd *cobra.Command, args []string) {
		serv, err := auth.NewAuthService(constants.AppID, constants.AppVersion, constants.AppSecret, "")
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(args) == 0 {
			fmt.Println("create: missing destination file path")
			return
		}
		newDir := userfile.NewFormattedPath(utils.GetCurrentOpDir(args, 0)).GetPath()
		sp := print.Spinner(os.Stdout, "Create Share [%s] ...", newDir)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		result, err := pubFileShare.NewPubFileShareClient(serv.GetGrpcConnection()).Create(ctx, &pubFileShare.FileShare{
			// Parent: &pubUserFile.File{Path: currentDir},
			PathList: []string{newDir},
		})
		sp(err == nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		data, _ := json.Marshal(result)
		fmt.Println(string(data))
	},
}

func init() {
	ShareCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
