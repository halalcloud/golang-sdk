/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package share

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	pubFileShare "github.com/city404/v6-public-rpc-proto/go/v6/fileshare"
	"github.com/halalcloud/golang-sdk/auth"
	"github.com/halalcloud/golang-sdk/constants"
	"github.com/spf13/cobra"
)

// getCmd represents the info command
var getCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"d"},
	Short:   "delete share",
	Long:    `Delete share info from the HalalCloud API.`,
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
		fmt.Printf("get task %s\n", deleteArray)

		serv, err := auth.NewAuthService(constants.AppID, constants.AppVersion, constants.AppSecret, "")
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, v := range deleteArray {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()
			result, err := pubFileShare.NewPubFileShareClient(serv.GetGrpcConnection()).Get(ctx, &pubFileShare.FileShare{
				Identity: v,
			})
			if err != nil {
				fmt.Println(err)
				return
			}
			data, _ := json.Marshal(result)
			fmt.Println(string(data))
		}
	},
}

func init() {
	ShareCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
