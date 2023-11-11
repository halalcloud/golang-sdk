package main

import (
	"context"

	common "github.com/city404/v6-public-rpc-proto/go/v6/common"
	pbPublicUser "github.com/city404/v6-public-rpc-proto/go/v6/user"
	pubUserFile "github.com/city404/v6-public-rpc-proto/go/v6/userfile"
	"github.com/halalcloud/golang-sdk/auth"
)

func main() {
	svc, err := auth.NewAuthServiceWithSimpleLogin("devDebugger/1.0", "1.0.0", "Nkx3Y2xvZ2luLmNu", "user", "password")
	if err != nil {
		panic(err)
	}
	// println(svc.GetAccessToken())

	grpcConnection := svc.GetGrpcConnection()
	defer grpcConnection.Close()

	userService := pbPublicUser.NewPubUserClient(grpcConnection)
	user, err := userService.Get(context.Background(), &pbPublicUser.User{})
	if err != nil {
		panic(err)
	}
	println(user.String())
	pubUserFileService := pubUserFile.NewPubUserFileClient(grpcConnection)
	userFile, err := pubUserFileService.Get(context.Background(), &pubUserFile.File{})
	if err != nil {
		panic(err)
	}
	println(userFile.String())
	pubUserList, err := pubUserFileService.List(context.Background(), &pubUserFile.FileListRequest{
		Parent: &pubUserFile.File{},
		ListInfo: &common.ScanListRequest{
			Limit: 5,
		},
	})
	if err != nil {
		panic(err)
	}
	println(pubUserList.String())
}
