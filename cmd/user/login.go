/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package user

import (
	"os"

	"github.com/halalcloud/golang-sdk/auth"
	"github.com/halalcloud/golang-sdk/constants"
	"github.com/halalcloud/golang-sdk/pkg/print"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//svc, err := auth.NewAuthServiceWithSimpleLogin("devDebugger/1.0", "1.0.0", "Nkx3Y2xvZ2luLmNu", "色情文件鉴定", "yipaihuyan")
		// if err != nil {
		//	panic(err)
		//}
		// println(svc.GetAccessToken())

		//grpcConnection := svc.GetGrpcConnection()
		//defer grpcConnection.Close()

		//userService := pbPublicUser.NewPubUserClient(grpcConnection)
		//user, err := userService.Get(context.Background(), &pbPublicUser.User{})
		//if err != nil {
		//	panic(err)
		//}
		//println(user.String())
		_, err := auth.NewAuthServiceWithOauth(os.Stdout, constants.AppID, constants.AppVersion, constants.AppSecret)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, "Login failed, %s", err.Error())
			return
		}
		// print.SuccessStatusEvent(os.Stdout, "Login success, %s", svc.GetAccessToken())
	},
}

func init() {
	UserCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
