/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package user

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	pbPublicUser "github.com/city404/v6-public-rpc-proto/go/v6/user"
	"github.com/halalcloud/golang-sdk/auth"
	"github.com/halalcloud/golang-sdk/constants"
	"github.com/spf13/cobra"
)

// iconCmd represents the login command
var iconCmd = &cobra.Command{
	Use:     "icon",
	Aliases: []string{"change-icon"},
	Short:   "A brief description of your command",
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
		newImage := args[0]
		f, err := os.Open(newImage)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		b, err := io.ReadAll(f)
		if err != nil {
			fmt.Println(err)
			return
		}
		b64 := base64.StdEncoding.EncodeToString(b)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		result, err := pbPublicUser.NewPubUserClient(serv.GetGrpcConnection()).Update(ctx, &pbPublicUser.User{
			Icon: b64,
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
	UserCmd.AddCommand(iconCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
