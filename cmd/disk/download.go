/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package disk

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	pubUserFile "github.com/city404/v6-public-rpc-proto/go/v6/userfile"
	"github.com/halalcloud/golang-sdk/auth"
	"github.com/halalcloud/golang-sdk/constants"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// "github.com/halalcloud/golang-sdk/pkg/downloader"
	"github.com/halalcloud/golang-sdk/pkg/print"
	"github.com/halalcloud/golang-sdk/utils"
	"github.com/ipfs/go-cid"
	"github.com/spf13/cobra"
	"github.com/zzzhr1990/go-common-entity/userfile"
)

// DownloadCmd represents the info command
var DownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a file",
	Long: `Download a file from the HalalCloud API.

Display Disk Usage, Quota.`,
	Run: func(cmd *cobra.Command, args []string) {
		serv, err := auth.NewAuthService(constants.AppID, constants.AppVersion, constants.AppSecret, "")
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(args) == 0 {
			fmt.Println("mkdir: missing operand")
			return
		}
		id, _ := cmd.Flags().GetString("id")

		newPath := userfile.NewFormattedPath(utils.GetCurrentOpDir(args, 0)).GetPath()
		if len(id) > 0 {
			newPath = "{id:{" + id + "}}"
		}

		sp := print.Spinner(os.Stdout, "Download Path [%s] ...", newPath)
		if len(id) > 0 {
			newPath = ""
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		client := pubUserFile.NewPubUserFileClient(serv.GetGrpcConnection())
		result, err := client.ParseFileSlice(ctx, &pubUserFile.File{

			// Parent: &pubUserFile.File{Path: currentDir},
			Path:     newPath,
			Identity: id,
		})
		if err != nil {
			sp(false)
			status, ok := status.FromError(err)
			if ok {
				if status.Code() == codes.NotFound {
					print.FailureStatusEvent(os.Stdout, "File [%s -> %s] not found, back to root.", newPath, id)
					return
				}
			}
			fmt.Println(err)
			return
		}
		fmt.Printf("source SHA1: %s, blocks: %d, file size: %d\n", result.Sha1, len(result.RawNodes), result.FileSize)
		fileAddrs := []*pubUserFile.SliceDownloadInfo{}
		batchRequest := []string{}
		for _, slice := range result.RawNodes {
			batchRequest = append(batchRequest, slice)
			if len(batchRequest) >= 200 {
				sliceAddress, err := client.GetSliceDownloadAddress(ctx, &pubUserFile.SliceDownloadAddressRequest{
					Identity: batchRequest,
				})
				if err != nil {
					sp(false)
					fmt.Println(err)
					return
				}
				fileAddrs = append(fileAddrs, sliceAddress.Addresses...)
				batchRequest = []string{}
			}
		}
		if len(batchRequest) > 0 {
			sliceAddress, err := client.GetSliceDownloadAddress(ctx, &pubUserFile.SliceDownloadAddressRequest{
				Identity: batchRequest,
			})
			if err != nil {
				sp(false)
				fmt.Println(err)
				return
			}
			fileAddrs = append(fileAddrs, sliceAddress.Addresses...)
		}
		sp(true)
		sha := sha1.New()

		fileSize := int64(0)
		for i, addr := range fileAddrs {

			displaySp := i%10 == 0
			var sp func(result print.Result) = nil
			if displaySp {
				sp = print.Spinner(os.Stdout, "Download Path [%d/%d] ...", i+1, len(fileAddrs))
			}
			dataBytes, err := tryAndGetRawFiles(addr)
			if err != nil {
				if sp != nil {
					sp(false)
				}
				fmt.Println(err)
				return
			}
			sha.Write(dataBytes)
			fileSize += int64(len(dataBytes))
			if sp != nil {
				sp(true)
			}
		}
		shaSum := sha.Sum(nil)
		fmt.Printf("SHA1: %x, file size: %d\n", shaSum, fileSize)
		// print.SuccessStatusEvent(os.Stdout, "File [%s -> %s] pieces: %s", newPath, id, string(dataJson))

		//ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		// defer cancel()
		// userFileClient := pubUserFile.NewPubUserFileClient(serv.GetGrpcConnection())
		//userFileClient.Get()
		//sliceDownloader := downloader.NewSliceDownloader(userFileClient)
		//err = sliceDownloader.StartDownload(newPath)
		//if err != nil {
		//	sp(false)
		//	print.FailureStatusEvent(os.Stdout, "Download [%s] failed, error: %s", newPath, err.Error())
		//	return
		//}
		// sp(true)
		print.SuccessStatusEvent(os.Stdout, "Download [%s] completed", newPath)
	},
}

func tryAndGetRawFiles(addr *pubUserFile.SliceDownloadInfo) ([]byte, error) {
	tryTimes := 0
	for {
		tryTimes++
		dataBytes, err := getRawFiles(addr)
		if err != nil {
			if tryTimes > 3 {
				return nil, err
			}
			continue
		}
		return dataBytes, nil
	}
}

func getRawFiles(addr *pubUserFile.SliceDownloadInfo) ([]byte, error) {
	client := http.Client{
		Timeout: time.Duration(60 * time.Second), // Set timeout to 5 seconds
	}
	resp, err := client.Get(addr.DownloadAddress)
	if err != nil {

		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s, body: %s", resp.Status, body)
	}

	sourceCid, err := cid.Decode(addr.Identity)
	if err != nil {
		return nil, err
	}
	checkCid, err := sourceCid.Prefix().Sum(body)
	if err != nil {
		return nil, err
	}
	if !checkCid.Equals(sourceCid) {
		return nil, fmt.Errorf("bad cid: %s, body: %s", checkCid.String(), body)
	}

	return body, nil

	// Process response body

	// Do something with the response body
}

func init() {
	DiskCmd.AddCommand(DownloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// DownloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// DownloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	DownloadCmd.Flags().StringP("id", "I", "", "rm by id")
}
