/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package disk

import (
	"fmt"
	"os"

	pubUserFile "github.com/city404/v6-public-rpc-proto/go/v6/userfile"
	"github.com/halalcloud/golang-sdk/auth"
	"github.com/halalcloud/golang-sdk/constants"
	"github.com/halalcloud/golang-sdk/pkg/downloader"
	"github.com/halalcloud/golang-sdk/pkg/print"
	"github.com/halalcloud/golang-sdk/utils"
	"github.com/spf13/cobra"
	"github.com/zzzhr1990/go-common-entity/userfile"
)

// DownloadCmd represents the info command
var DownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a file",
	Long: `Download a file from the HalalCloud API.

Display Disk Usage, Quota.`,
	Run: func(cmd *cobra.Command, args []string) {
		serv, err := auth.NewAuthService(constants.AppID, constants.AppVersion, constants.AppSecret, "")
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(args) == 0 {
			fmt.Println("mkdir: missing operand")
			return
		}
		force, ok := cmd.Flags().GetBool("force")
		if ok == nil && force {
			print.WarningStatusEvent(os.Stdout, "Force remove, instead of put items into trash.")
		} else {
			cmd.Run(trashCmd, args)
			return
		}
		id, _ := cmd.Flags().GetString("id")

		newPath := userfile.NewFormattedPath(utils.GetCurrentOpDir(args, 0)).GetPath()
		if len(id) > 0 {
			newPath = "{id:{" + id + "}}"
		}

		sp := print.Spinner(os.Stdout, "Remove Path [%s] ...", newPath)
		if len(id) > 0 {
			newPath = ""
		}
		//ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		// defer cancel()
		userFileClient := pubUserFile.NewPubUserFileClient(serv.GetGrpcConnection())
		sliceDownloader := downloader.NewSliceDownloader(userFileClient)
		err = sliceDownloader.StartDownload(id)
		if err != nil {
			sp(false)
			print.FailureStatusEvent(os.Stdout, "Download [%s] failed", newPath)
			return
		}
		sp(true)
		print.SuccessStatusEvent(os.Stdout, "Download [%s] completed", newPath)
	},
}

func init() {
	DiskCmd.AddCommand(DownloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// DownloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// DownloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
