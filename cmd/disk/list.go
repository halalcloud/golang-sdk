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
	"github.com/dustin/go-humanize"
	"github.com/eiannone/keyboard"
	"github.com/halalcloud/golang-sdk/auth"
	"github.com/halalcloud/golang-sdk/constants"
	cPrint "github.com/halalcloud/golang-sdk/pkg/print"
	"github.com/halalcloud/golang-sdk/utils"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var ListCmd = &cobra.Command{
	Use:     "list",
	Short:   "list files in current directory",
	Long:    `list files in current directory`,
	Aliases: []string{"ls", "dir"},
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
		opDir := utils.GetCurrentDir()
		if len(args) > 0 {
			opDir = utils.GetCurrentOpDir(args, 0)
		}
		client := pubUserFile.NewPubUserFileClient(serv.GetGrpcConnection())
		for {
			timeStart := time.Now()
			sp := cPrint.Spinner(os.Stdout, "Listing Directory [%s] ...", opDir)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()
			result, err := client.List(ctx, &pubUserFile.FileListRequest{
				Parent: &pubUserFile.File{Path: opDir},
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
				printList(opDir, result)
			}
			cPrint.InfoStatusEvent(os.Stdout, "%d items, %s escaped.", len(result.Files), timeEscaped.String())
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
	DiskCmd.AddCommand(ListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func printList(_ string, result *pubUserFile.FileListResponse) {
	// data := [][]string{}

	// table := tablewriter.NewWriter(os.Stdout)
	// table.SetHeader([]string{"Name", "Type", "ID", "Size", "CreateTime"})
	for _, v := range result.Files {
		size := "-"
		if !v.Dir {
			size = humanize.IBytes(uint64(v.Size))
		}
		// print.Info(os.Stdout, "%s %s %s %s %s", v.Name, v.MimeType, v.Identity, size, time.UnixMilli(v.CreateTs).Format("2006-01-02 15:04:05"))
		// table.Append([]string{v.Name, v.MimeType, v.Identity, size, time.UnixMilli(v.CreateTs).Format("2006-01-02 15:04:05")})
		println(fmt.Sprintf("*************************\nName: %s\nMime: %s\nID: %s\nPath: %s\nSize: %s (v:%d)\nCreate Time:%s\n*************************\n",
			v.Name,
			v.MimeType,
			v.Identity,
			v.Path,
			size,
			v.Version,
			time.UnixMilli(v.UpdateTs).Format("2006-01-02 15:04:05")))
	}

	// hasMore := false
	// if result.ListInfo != nil && result.ListInfo.Token != "" {
	// 	hasMore = true
	//}

	//if hasMore {
	//	table.SetFooter([]string{"", "", "", "", "..."}) // Add Footer
	//}

	// table.EnableBorder(false)                             // Set Border to false
	//table.AppendBulk(data) // Add Bulk Data
	//table.Render()
}
