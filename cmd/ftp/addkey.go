/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package ftp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pbSftpConfig "github.com/city404/v6-public-rpc-proto/go/v6/sftpconfig"
	"github.com/halalcloud/golang-sdk/auth"
	"github.com/halalcloud/golang-sdk/constants"
	"github.com/spf13/cobra"
)

// addKeyCmd represents the login command
var addKeyCmd = &cobra.Command{
	Use:   "addkey",
	Short: "Add ssh key to the server",
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

		if len(args) < 1 {
			fmt.Println("Please provide the ssh key")
			return
		}

		sshKey := args[0]
		if len(sshKey) < 1 {
			fmt.Println("Please provide the ssh key")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		client := pbSftpConfig.NewPubSftpConfigClient(serv.GetGrpcConnection())
		result, err := client.Get(ctx, &pbSftpConfig.SftpConfig{})
		if err != nil {
			fmt.Println(err)
			return
		}
		sourceKey := result.SshKey
		if len(sourceKey) > 0 {
			sshKey = sourceKey + "\n" + sshKey
		}
		data, err := client.Update(ctx, &pbSftpConfig.SftpConfig{SshKey: sshKey})
		if err != nil {
			fmt.Println(err)
			return
		}
		dataJson, _ := json.Marshal(data)
		fmt.Println(string(dataJson))
	},
}

func init() {
	FtpCmd.AddCommand(addKeyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
