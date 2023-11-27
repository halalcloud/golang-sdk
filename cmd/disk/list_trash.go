/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package disk

import (
	"context"
	"fmt"
	"os"
	"time"

	common "github.com/city404/v6-public-rpc-proto/go/v6/common"
	pubUserFile "github.com/city404/v6-public-rpc-proto/go/v6/userfile"
	"github.com/eiannone/keyboard"
	"github.com/halalcloud/golang-sdk/auth"
	"github.com/halalcloud/golang-sdk/constants"
	"github.com/halalcloud/golang-sdk/pkg/print"
	"github.com/spf13/cobra"
)

// listTrashCmd represents the listTrash command
var listTrashCmd = &cobra.Command{
	Use:   "list-trash",
	Short: "List trash items",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		serv, err := auth.NewAuthService(constants.AppID, constants.AppVersion, constants.AppSecret, "")
		if err != nil {
			fmt.Println(err)
			return
		}
		limit := int64(10)
		token := ""
		keyEnable := true
		if err := keyboard.Open(); err != nil {
			keyEnable = false
		}
		defer func() {
			_ = keyboard.Close()
		}()

		client := pubUserFile.NewPubUserFileClient(serv.GetGrpcConnection())
		for {
			timeStart := time.Now()
			sp := print.Spinner(os.Stdout, "Listing Trash  ...")
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()
			result, err := client.ListTrash(ctx, &pubUserFile.FileListRequest{
				Parent: &pubUserFile.File{},
				ListInfo: &common.ScanListRequest{
					Limit: limit,
					Token: token,
				},
			})
			if err != nil {
				sp(false)
				fmt.Println(err)
				return
			}
			sp(true)
			timeEscaped := time.Since(timeStart)

			if result.Files != nil && len(result.Files) > 0 {
				printList("-", result)
			}
			print.InfoStatusEvent(os.Stdout, "%d items, %s escaped.", len(result.Files), timeEscaped.String())
			if result.ListInfo == nil || result.ListInfo.Token == "" {
				break
			}
			token = result.ListInfo.Token
			if keyEnable {
				fmt.Println("More Results, Press any key to continue..., Press ESC/Ctrl+C to exit.")
				_, key, err := keyboard.GetSingleKey()
				if err != nil {
					continue
				} else if key == keyboard.KeyEsc || key == keyboard.KeyCtrlC {
					break
				}
			}
		}
	},
}

func init() {
	DiskCmd.AddCommand(listTrashCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listTrashCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listTrashCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
