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

// MvCmd represents the  command
var RnCmd = &cobra.Command{
	Use:     "rn",
	Short:   "rename file or directory",
	Aliases: []string{"rename"},
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		currentDir := utils.GetCurrentDir()
		serv, err := auth.NewAuthService(constants.AppID, constants.AppVersion, constants.AppSecret, "")
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(args) == 0 {
			fmt.Println("rename: missing operand")
			return
		}
		if len(args) == 1 {
			fmt.Println("rename: missing destination file operand after '" + args[0] + "'")
			return
		}

		id, _ := cmd.Flags().GetString("id")

		newPath := userfile.NewFormattedPath(utils.GetCurrentOpDir(args, 0)).GetPath()
		if len(id) > 0 {
			newPath = "{id:{" + id + "}}"
		}

		sp := print.Spinner(os.Stdout, "rename Path [%s] -> ...", newPath)
		if len(id) > 0 {
			newPath = ""
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		result, err := pubUserFile.NewPubUserFileClient(serv.GetGrpcConnection()).Rename(ctx, &pubUserFile.File{

			// Parent: &pubUserFile.File{Path: currentDir},
			Path:     newPath,
			Identity: id,
			Name:     args[1],
		},

		// Dest: dest,
		)
		if err != nil {
			sp(false)
			status, ok := status.FromError(err)
			if ok {
				if status.Code() == codes.NotFound {
					print.FailureStatusEvent(os.Stdout, "Destination [%s] not found. id=%s, path=%s", currentDir, newPath, id)
					// utils.SetCurrentDir("/")
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
		print.SuccessStatusEvent(os.Stdout, "Rename operation [%s] created. status: %d", result.Task, result.Status)
	},
}

func init() {
	DiskCmd.AddCommand(RnCmd)

	RnCmd.Flags().StringP("id", "I", "", "rename by id")
}
