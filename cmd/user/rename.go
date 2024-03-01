/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package user

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pbPublicUser "github.com/city404/v6-public-rpc-proto/go/v6/user"
	"github.com/halalcloud/golang-sdk/auth"
	"github.com/halalcloud/golang-sdk/constants"
	"github.com/spf13/cobra"
)

// getCmd represents the login command
var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "A brief description of your command",
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
			fmt.Println("rename: missing operand")
			return
		}
		newName := args[0]
		if len(newName) < 2 {
			fmt.Println("rename: name too short")
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		result, err := pbPublicUser.NewPubUserClient(serv.GetGrpcConnection()).Update(ctx, &pbPublicUser.User{
			Name: newName,
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
	UserCmd.AddCommand(renameCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
