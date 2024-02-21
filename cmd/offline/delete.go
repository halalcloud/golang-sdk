/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package offline

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	pubUserOffline "github.com/city404/v6-public-rpc-proto/go/v6/offline"
	"github.com/halalcloud/golang-sdk/auth"
	"github.com/halalcloud/golang-sdk/constants"
	"github.com/spf13/cobra"
)

// deleteCmd represents the info command
var deleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"del"},
	Short:   "delete url info to offline download",
	Long:    `Delete offline download info from the HalalCloud API.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
		deleteArray := []string{}
		for _, v := range args {
			deleteArray = append(deleteArray, strings.Split(v, ",")...)

			// deleteArray = append(deleteArray, v)
		}
		fmt.Printf("delete task %s\n", deleteArray)

		serv, err := auth.NewAuthService(constants.AppID, constants.AppVersion, constants.AppSecret, "")
		if err != nil {
			fmt.Println(err)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		result, err := pubUserOffline.NewPubOfflineTaskClient(serv.GetGrpcConnection()).Delete(ctx, &pubUserOffline.OfflineTaskDeleteRequest{
			Identity: deleteArray,
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
	OfflineCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
