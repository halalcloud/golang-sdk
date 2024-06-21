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
	"path/filepath"
	"time"

	pubUserFile "github.com/city404/v6-public-rpc-proto/go/v6/userfile"
	"github.com/halalcloud/golang-sdk/auth"
	"github.com/halalcloud/golang-sdk/constants"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// "github.com/halalcloud/golang-sdk/pkg/downloader"
	"github.com/halalcloud/golang-sdk/pkg/downloader"
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
		id, _ := cmd.Flags().GetString("id")
		if len(args) == 0 && len(id) == 0 {
			fmt.Println("download: missing operand")

			return
		}

		newPath := ""
		if len(id) > 0 {
			newPath = "{id:{" + id + "}}"
		} else {
			newPath = userfile.NewFormattedPath(utils.GetCurrentOpDir(args, 0)).GetPath()
		}
		client := pubUserFile.NewPubUserFileClient(serv.GetGrpcConnection())
		dirname, err := os.UserHomeDir()
		if err != nil {
			print.FailureStatusEvent(os.Stdout, "Get home dir failed, error: %s", err.Error())
		}
		downloadsDir := filepath.Join(dirname, "Downloads")
		print.InfoStatusEvent(os.Stdout, "Download Path [%s] ...", downloadsDir)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		sliceDownloader := downloader.NewSliceDownloader(client, ctx, "")
		err = sliceDownloader.Start()
		if err != nil {
			fmt.Println(err)
			return
		}
		defer sliceDownloader.Stop()

		sp := print.Spinner(os.Stdout, "Download Path [%s] ...", newPath)
		if len(id) > 0 {
			newPath = ""
		}

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
					print.FailureStatusEvent(os.Stdout, "1File [%s -> %s] not found, back to root.", newPath, id)
					return
				}
			}
			fmt.Printf("slice parse error: %v\r", err)
			return
		}
		filePath := filepath.Join(findDownloadDir(), result.Name)
		fmt.Printf("source SHA1: %s, blocks: %d, file size: %d, file: %s\n", result.Sha1, len(result.RawNodes), result.FileSize, filePath)
		// fileName := result.
		fileAddrs := []*pubUserFile.SliceDownloadInfo{}
		batchRequest := []string{}

		for _, slice := range result.RawNodes {
			batchRequest = append(batchRequest, slice)
			if len(batchRequest) >= 200 {
				sliceAddress, err := client.GetSliceDownloadAddress(ctx, &pubUserFile.SliceDownloadAddressRequest{
					Identity: batchRequest,
					Version:  1,
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
				Version:  1,
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
		//filePath := filepath.Join(utils.GetCurrentDir(), result)

		fileDownload, err := os.Create(filePath)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer fileDownload.Close()
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
			_, err = fileDownload.Write(dataBytes)
			if err != nil {
				if sp != nil {
					sp(false)
				}
				fmt.Println(err)
				return
			}
			fileSize += int64(len(dataBytes))
			if sp != nil {
				sp(true)
			}
		}
		shaSum := sha.Sum(nil)
		fmt.Printf("SHA1: %x, file size: %d -> %s\n", shaSum, fileSize, filePath)
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

	if addr.Encrypt > 0 {
		cd := uint8(addr.Encrypt)
		for idx := 0; idx < len(body); idx++ {
			body[idx] = body[idx] ^ cd
		}
	}

	if addr.StoreType != 10 {

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

var (
	DownloadDirNames []string = []string{"Downloads", "downloads", "download", "Downloads", "etc..."}
)

func findDownloadDir() string {
	// var downloadDir string
	var DownloadDirNames []string = []string{"Downloads", "downloads", "download", "Downloads", "etc..."}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return os.TempDir()
	}

	for _, ddn := range DownloadDirNames {
		var dir = filepath.Join(homeDir, ddn)

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			// fmt.Println(dir, "does not exist")
		} else {
			// fmt.Println("The provided directory named", dir, "exists")
			return dir
		}
	}
	return os.TempDir()
}
