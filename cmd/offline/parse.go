/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package offline

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pubUserOffline "github.com/city404/v6-public-rpc-proto/go/v6/offline"
	"github.com/halalcloud/golang-sdk/auth"
	"github.com/halalcloud/golang-sdk/constants"
	"github.com/spf13/cobra"
)

// parseCmd represents the info command
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse url info",
	Long: `Get offline download info from the HalalCloud API.

Display Disk Usage, Quota.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("parse called")
		if len(args) == 0 {
			cmd.Help()
			return
		}
		fmt.Println("parse called with", args[0])
		serv, err := auth.NewAuthService(constants.AppID, constants.AppVersion, constants.AppSecret, "")
		if err != nil {
			fmt.Println(err)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		result, err := pubUserOffline.NewPubOfflineTaskClient(serv.GetGrpcConnection()).Parse(ctx, &pubUserOffline.TaskParseRequest{
			Url: args[0],
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		data, _ := json.Marshal(result)
		fmt.Println(string(data))
	},
}

func init() {
	OfflineCmd.AddCommand(parseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
