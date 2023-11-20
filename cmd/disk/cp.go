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

// cpCmd represents the cp command
var CpCmd = &cobra.Command{
	Use:     "cp",
	Short:   "copy file or directory",
	Aliases: []string{"copy"},
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
			fmt.Println("mkdir: missing operand")
			return
		}
		if len(args) == 1 {
			fmt.Println("mkdir: missing destination file operand after '" + args[0] + "'")
			return
		}

		id, _ := cmd.Flags().GetString("id")
		identityByPath := false

		newPath := userfile.NewFormattedPath(utils.GetCurrentOpDir(args, 0)).GetPath()
		if len(id) > 0 {
			newPath = "{id:{" + id + "}}"
			identityByPath = true
		}

		sp := print.Spinner(os.Stdout, "Copy Path [%s] ...", newPath)
		if len(id) > 0 {
			newPath = ""
		}
		dest := &pubUserFile.File{}
		if identityByPath {
			dest.Identity = args[1]
		} else {
			dest.Path = userfile.NewFormattedPath(utils.GetCurrentOpDir(args, 1)).GetPath()
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		result, err := pubUserFile.NewPubUserFileClient(serv.GetGrpcConnection()).Copy(ctx, &pubUserFile.BatchOperationRequest{
			Source: []*pubUserFile.File{
				{
					// Parent: &pubUserFile.File{Path: currentDir},
					Path:     newPath,
					Identity: id,
				},
			},
			Dest: dest,
		})
		if err != nil {
			sp(false)
			status, ok := status.FromError(err)
			if ok {
				if status.Code() == codes.NotFound {
					print.FailureStatusEvent(os.Stdout, "Destination [%s] not found. id=%s, path=%s", currentDir, dest.Identity, dest.Path)
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
		print.SuccessStatusEvent(os.Stdout, "Remove operation [%s] created.", result.Task)
	},
}

func init() {
	DiskCmd.AddCommand(CpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	CpCmd.Flags().StringP("id", "I", "", "copy by id")
}
