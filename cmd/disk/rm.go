/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package disk

import (
	"context"
	"fmt"
	"os"
	"time"

	pubUserFile "github.com/city404/v6-public-rpc-proto/go/v6/userfile"
	"github.com/halalcloud/golang-sdk/auth"
	"github.com/halalcloud/golang-sdk/constants"
	"github.com/halalcloud/golang-sdk/pkg/print"
	"github.com/halalcloud/golang-sdk/utils"
	"github.com/spf13/cobra"
	"github.com/zzzhr1990/go-common-entity/userfile"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// rmCmd represents the rm command
var RmCmd = &cobra.Command{
	Use:   "rm",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Aliases: []string{"del", "remove"},
	Run: func(cmd *cobra.Command, args []string) {
		currentDir := utils.GetCurrentDir()
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
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		result, err := pubUserFile.NewPubUserFileClient(serv.GetGrpcConnection()).Delete(ctx, &pubUserFile.BatchOperationRequest{
			Source: []*pubUserFile.File{
				{
					// Parent: &pubUserFile.File{Path: currentDir},
					Path:     newPath,
					Identity: id,
				},
			},
		})
		if err != nil {
			sp(false)
			status, ok := status.FromError(err)
			if ok {
				if status.Code() == codes.NotFound {
					print.FailureStatusEvent(os.Stdout, "Directory [%s] not found, back to root.", currentDir)
					utils.SetCurrentDir("/")
					return
				}
				if status.Code() == codes.FailedPrecondition {
					// fmt.Printf("Directory [%s] is locking, back to root.\n", currentDir)
					print.FailureStatusEvent(os.Stdout, "Directory [%s] is locking, please wait background progress completed.", currentDir)
					return
				}
			}
			fmt.Println(err)
			return
		}
		sp(true)
		print.SuccessStatusEvent(os.Stdout, "Remove operation [%s] created.", result.Task)
	},
}

func init() {
	DiskCmd.AddCommand(RmCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	RmCmd.Flags().BoolP("force", "F", false, "Force remove, instead of put items into trash.")
	RmCmd.Flags().StringP("id", "I", "", "rm by id")
}
