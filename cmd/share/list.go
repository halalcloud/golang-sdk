/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package share

import (
	"context"
	"os"
	"strconv"

	"github.com/dustin/go-humanize"
	"github.com/halalcloud/golang-sdk/pkg/print"
	"github.com/olekukonko/tablewriter"

	"fmt"
	"time"

	common "github.com/city404/v6-public-rpc-proto/go/v6/common"
	pubFileShare "github.com/city404/v6-public-rpc-proto/go/v6/fileshare"
	"github.com/eiannone/keyboard"
	"github.com/halalcloud/golang-sdk/auth"
	"github.com/halalcloud/golang-sdk/constants"
	"github.com/spf13/cobra"
)

// listCmd represents the info command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "list user share",
	Long:    `List share.`,
	Run: func(cmd *cobra.Command, args []string) {

		serv, err := auth.NewAuthService(constants.AppID, constants.AppVersion, constants.AppSecret, "")
		if err != nil {
			fmt.Println(err)
			return
		}

		client := pubFileShare.NewPubFileShareClient(serv.GetGrpcConnection())
		limit := int64(10)
		token := ""
		keyEnable := true
		if err := keyboard.Open(); err != nil {
			keyEnable = false
		}
		defer func() {
			_ = keyboard.Close()
		}()
		for {
			timeStart := time.Now()
			sp := print.Spinner(os.Stdout, "Listing Offline Download List ...")
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()
			result, err := client.List(ctx, &pubFileShare.FileShareListRequest{
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

			if result.Shares != nil && len(result.Shares) > 0 {
				printList(result)
			}
			print.InfoStatusEvent(os.Stdout, "%d items, %s escaped.", len(result.Shares), timeEscaped.String())
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
	ShareCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func printList(result *pubFileShare.FileShareListResponse) {
	data := [][]string{}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "ADDR", "ID", "Size", "CreateTime"})
	for _, v := range result.Shares {
		size := "-"

		size = humanize.IBytes(uint64(v.FileSize))

		table.Append([]string{v.Name, strconv.Itoa(len(v.PathList)), v.Identity, size, time.UnixMilli(v.CreateTs).Format("2006-01-02 15:04:05")})
	}

	hasMore := false
	if result.ListInfo != nil && result.ListInfo.Token != "" {
		hasMore = true
	}

	if hasMore {
		table.SetFooter([]string{"", "", "", "", "..."}) // Add Footer
	}

	// table.EnableBorder(false)                             // Set Border to false
	table.AppendBulk(data) // Add Bulk Data
	table.Render()
}
