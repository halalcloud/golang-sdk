/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package disk

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	bauth "github.com/baidubce/bce-sdk-go/auth"
	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/bos"
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

// UploadCmd represents the mkdir command
var UploadCmd = &cobra.Command{
	Use:     "upload",
	Short:   "A brief description of your command",
	Aliases: []string{"up"},
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
			fmt.Println("create: missing operand")
			return
		}

		// get last arg, and check if it is a file
		lastArg := args[len(args)-1]
		if !utils.IsFile(lastArg) {
			fmt.Println("create: last operand " + lastArg + " must be a file")
			return
		}
		path, _ := filepath.Abs(lastArg)
		fileInfo, err := os.Stat(path)
		if err != nil {
			fmt.Println(err)
			return
		}

		log.Printf("file: %s, file path: %s", fileInfo.Name(), path)

		newDir := userfile.NewFormattedPath(utils.GetCurrentDir()).GetPath()
		newDir = strings.TrimSuffix(newDir, "/") + "/" + fileInfo.Name()
		sp := print.Spinner(os.Stdout, "Upload File [%s] ...", newDir)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		result, err := pubUserFile.NewPubUserFileClient(serv.GetGrpcConnection()).CreateUploadToken(ctx, &pubUserFile.File{
			// Parent: &pubUserFile.File{Path: currentDir},
			Path: newDir,
			//ContentIdentity: args[1],
		})
		if err != nil {
			sp(false)
			status, ok := status.FromError(err)
			if ok {
				if status.Code() == codes.NotFound {
					fmt.Printf("Directory [%s] not found, back to root.\n", currentDir)
					utils.SetCurrentDir("/")
					return
				}
			}
			fmt.Println(err)
			return
		}
		sp(true)

		print.InfoStatusEvent(os.Stdout, "Starting upload [%s].", path)
		// -^-^-
		clientConfig := bos.BosClientConfiguration{
			Ak:               result.AccessKey,
			Sk:               result.SecretKey,
			Endpoint:         result.Endpoint,
			RedirectDisabled: false,
			//SessionToken:     result.SessionToken,
		}

		// 初始化一个BosClient
		bosClient, err := bos.NewClientWithConfig(&clientConfig)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, "failed to create bos client: %v", err)
			return
		}
		stsCredential, err := bauth.NewSessionBceCredentials(
			result.AccessKey,
			result.SecretKey,
			result.Token)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, "failed to create sts credential: %v", err)
			return
		}
		bosClient.Config.Credentials = stsCredential
		// return bosClient, nil
		body, err := bce.NewBodyFromFile(path)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, "failed to read file: %v", err)
			return
		}
		postResult, err := bosClient.PutObject(result.Bucket, result.Key, body, nil)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, "failed to upload file: %v ===> %s/%s", err, clientConfig.Ak, clientConfig.Sk)
			return
		}
		print.SuccessStatusEvent(os.Stdout, "File [%s] uploaded. etag: %s, key: %s", fileInfo.Name(), postResult, result.Key)
	},
}

func init() {
	DiskCmd.AddCommand(UploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mkdirCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mkdirCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
