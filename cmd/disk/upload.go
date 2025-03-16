/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package disk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	pubUserFile "github.com/city404/v6-public-rpc-proto/go/v6/userfile"
	"github.com/halalcloud/golang-sdk/auth"
	"github.com/halalcloud/golang-sdk/constants"
	"github.com/halalcloud/golang-sdk/utils"
	"github.com/ipfs/go-cid"
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
		serv, err := auth.NewAuthService(constants.AppID, constants.AppVersion, constants.AppSecret, "")
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(args) == 0 {
			fmt.Println("create: missing operand")
			return
		}
		serv.GetGrpcConnection()

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
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		result, err := pubUserFile.NewPubUserFileClient(serv.GetGrpcConnection()).CreateUploadTask(ctx, &pubUserFile.File{
			// Parent: &pubUserFile.File{Path: currentDir},
			Path: newDir,
			//ContentIdentity: args[1],
			Size: fileInfo.Size(),
		})
		if err != nil {
			status, ok := status.FromError(err)
			if ok {
				if status.Code() == codes.NotFound {
					fmt.Printf("Directory [%s] not found, back to root.\n", utils.GetCurrentDir())
					utils.SetCurrentDir("/")
					return
				}
			}
			fmt.Println(err)
			return
		}
		if result.Created {
			fmt.Printf("Upload file created\n")
			return
		}
		fmt.Printf("Upload task started, block size: %d -> %s\n", result.BlockSize, result.Task)
		bufferSize := int(result.BlockSize)
		buffer := make([]byte, bufferSize)
		slicesList := make([]string, 0)
		codec := uint64(cid.Raw)
		if result.BlockCodec > 0 {
			codec = uint64(result.BlockCodec)
		}
		mhType := uint64(0x12)
		if result.BlockHashType > 0 {
			mhType = uint64(result.BlockHashType)
		}
		prefix := cid.Prefix{
			Codec:    codec,
			MhLength: -1,
			MhType:   mhType,
			Version:  1,
		}
		// read file
		fi, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer fi.Close()
		for {
			n, err := fi.Read(buffer)
			if n > 0 {
				data := buffer[:n]
				uploadCid, err := postFileSlice(data, result.Task, result.UploadAddress, prefix)
				if err != nil {
					fmt.Println(err)
					return
				}
				slicesList = append(slicesList, uploadCid.String())
			}
			if err == io.EOF || n == 0 {
				break
			}
		}
		newFile, err := makeFile(slicesList, result.Task, result.UploadAddress)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("File uploaded, cid: %s\n", newFile.ContentIdentity)

		// helpers.BlockSizeLimit = constants.IpfsDefaultBlockSize * 2
		/*
			sp := print.Spinner(os.Stdout, "Upload File [%s] ...", newDir)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()
			result, err := pubUserFile.NewPubUserFileClient(serv.GetGrpcConnection()).CreateUploadTask(ctx, &pubUserFile.File{
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
		*/

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
func makeFile(fileSlice []string, taskID string, uploadAddress string) (*pubUserFile.File, error) {
	accessUrl := uploadAddress + "/" + taskID
	getTimeOut := time.Minute * 2
	u, err := url.Parse(accessUrl)
	if err != nil {
		return nil, err
	}
	n, _ := json.Marshal(fileSlice)
	httpRequest := http.Request{
		Method: http.MethodPost,
		URL:    u,
		Header: map[string][]string{
			"Accept":         {"application/json"},
			"Content-Type":   {"application/octet-stream"},
			"Content-Length": {fmt.Sprintf("%d", len(fileSlice))},
		},
		Body: io.NopCloser(bytes.NewReader(n)),
	}
	httpClient := http.Client{
		Timeout: getTimeOut,
	}
	httpResponse, err := httpClient.Do(&httpRequest)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != http.StatusOK && httpResponse.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(httpResponse.Body)
		fmt.Println(string(b))
		return nil, fmt.Errorf("mk file slice failed, status code: %d", httpResponse.StatusCode)
	}
	b, _ := io.ReadAll(httpResponse.Body)
	var result *pubUserFile.File
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func postFileSlice(fileSlice []byte, taskID string, uploadAddress string, preix cid.Prefix) (cid.Cid, error) {
	// 1. sum file slice
	newCid, err := preix.Sum(fileSlice)
	if err != nil {
		return cid.Undef, err
	}
	// 2. post file slice
	sliceCidString := newCid.String()
	// /{taskID}/{sliceID}
	accessUrl := uploadAddress + "/" + taskID + "/" + sliceCidString
	getTimeOut := time.Second * 30
	// get {accessUrl} in {getTimeOut}
	u, err := url.Parse(accessUrl)
	if err != nil {
		return cid.Undef, err
	}
	// header: accept: application/json
	// header: content-type: application/octet-stream
	// header: content-length: {fileSlice.length}
	// header: x-content-cid: {sliceCidString}
	// header: x-task-id: {taskID}
	httpRequest := http.Request{
		Method: http.MethodGet,
		URL:    u,
		Header: map[string][]string{
			"Accept": {"application/json"},
		},
	}
	httpClient := http.Client{
		Timeout: getTimeOut,
	}
	httpResponse, err := httpClient.Do(&httpRequest)
	if err != nil {
		return cid.Undef, err
	}
	if httpResponse.StatusCode != http.StatusOK {
		return cid.Undef, fmt.Errorf("upload file slice failed, status code: %d", httpResponse.StatusCode)
	}
	var result bool
	b, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return cid.Undef, err
	}
	err = json.Unmarshal(b, &result)
	if err != nil {
		return cid.Undef, err
	}
	if result {
		return newCid, nil
	}

	httpRequest = http.Request{
		Method: http.MethodPost,
		URL:    u,
		Header: map[string][]string{
			"Accept":         {"application/json"},
			"Content-Type":   {"application/octet-stream"},
			"Content-Length": {fmt.Sprintf("%d", len(fileSlice))},
		},
		Body: io.NopCloser(bytes.NewReader(fileSlice)),
	}
	httpResponse, err = httpClient.Do(&httpRequest)
	if err != nil {
		return cid.Undef, err
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != http.StatusOK && httpResponse.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(httpResponse.Body)
		fmt.Println(string(b))
		return cid.Undef, fmt.Errorf("upload file slice failed, status code: %d", httpResponse.StatusCode)
	}
	//

	return newCid, nil
}
