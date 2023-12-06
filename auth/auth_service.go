package auth

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"

	pbPublicUser "github.com/city404/v6-public-rpc-proto/go/v6/user"
	"github.com/google/uuid"
	"github.com/halalcloud/golang-sdk/pkg/print"
	"github.com/halalcloud/golang-sdk/utils"
	"github.com/mdp/qrterminal/v3"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	appID                string
	appVersion           string
	appSecret            string
	refreshToken         string
	accessToken          string
	accessTokenExpiredAt int64
	grpcConnection       *grpc.ClientConn
	dopts                halalOptions
	configMutex          sync.RWMutex
}

func NewAuthServiceWithOauth(writer io.Writer, appID, appVersion, appSecret string, options ...HalalOption) (*AuthService, error) {
	svc := &AuthService{
		appID:        appID,
		appVersion:   appVersion,
		appSecret:    appSecret,
		refreshToken: "",
		dopts:        defaultOptions(),
	}
	for _, opt := range options {
		opt.apply(&svc.dopts)
	}

	grpcServer := "grpcuserapi.2dland.cn:443"
	dialContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	grpcOptions := svc.dopts.grpcOptions
	grpcOptions = append(grpcOptions, grpc.WithBlock(), grpc.WithAuthority("grpcuserapi.2dland.cn"), grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})), grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctxx := svc.signContext(method, ctx)
		err := invoker(ctxx, method, req, reply, cc, opts...) // invoking RPC method
		return err
	}))

	grpcConnection, err := grpc.DialContext(dialContext, grpcServer, grpcOptions...)
	if err != nil {
		return nil, err
	}
	defer grpcConnection.Close()
	userClient := pbPublicUser.NewPubUserClient(grpcConnection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stateString := uuid.New().String()
	// queryValues.Add("callback", oauthToken.Callback)
	oauthToken, err := userClient.CreateAuthToken(ctx, &pbPublicUser.LoginRequest{
		ReturnType: 2,
		State:      stateString,
		ReturnUrl:  "",
	})
	if err != nil {
		return nil, err
	}
	if len(oauthToken.State) < 1 {
		oauthToken.State = stateString
	}

	// &scope=nodata&state=DrcNxZqgLUb8nIYHAxbWSG82UaPjx2BYzE9G6jba&callback=nodata
	fmt.Fprintln(os.Stderr, "Plase open this url in your browser to login:\r\n\r\n"+oauthToken.Url+"\r\nOr san this QR Code：")

	config := qrterminal.Config{
		Level:     qrterminal.M,
		Writer:    os.Stdout,
		BlackChar: qrterminal.WHITE,
		WhiteChar: qrterminal.BLACK,
		QuietZone: 1,
	}
	qrterminal.GenerateWithConfig(oauthToken.Url, config)
	checkTimer := time.NewTicker(5 * time.Second)
	defer checkTimer.Stop()
	// catch interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	print.PendingStatusEvent(writer, "Waiting for login...[%s]", oauthToken.Url)
	for {
		select {
		case <-checkTimer.C:

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			checkLoginResponse, err := userClient.VerifyAuthToken(ctx, &pbPublicUser.LoginRequest{
				State:      oauthToken.State,
				Callback:   oauthToken.Callback,
				ReturnType: 2,
			})
			if err != nil {
				return nil, err
			}
			if checkLoginResponse.Status == 6 {
				login := checkLoginResponse.Login
				if login == nil {
					return nil, fmt.Errorf("login is nil")
				}
				if login.User != nil && len(login.Token.RefreshToken) > 0 {
					// checkLoginResponse = checkLoginResponse
					svc.OnAccessTokenRefreshed(login.Token.AccessToken, login.Token.AccessTokenExpireTs, login.Token.RefreshToken, login.Token.RefreshTokenExpireTs)
					newAuthService, err := NewAuthService(appID, appVersion, appSecret, login.Token.RefreshToken)
					if err != nil {
						return nil, err
					}
					return newAuthService, nil
					// break
				}
			}

			// reset timer
			checkTimer.Reset(5 * time.Second)
		case <-c:
			// interrupt signal
			return nil, fmt.Errorf("interrupted")
		}
	}

}

func NewAuthService(appID, appVersion, appSecret, refreshToken string, options ...HalalOption) (*AuthService, error) {

	svc := &AuthService{
		appID:        appID,
		appVersion:   appVersion,
		appSecret:    appSecret,
		refreshToken: refreshToken,
		dopts:        defaultOptions(),
	}

	if len(refreshToken) < 1 {
		refreshToken = viper.GetString("refresh_token")
	}

	for _, opt := range options {
		opt.apply(&svc.dopts)
	}

	grpcServer := "grpcuserapi.2dland.cn:443"
	dialContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	grpcOptions := svc.dopts.grpcOptions
	grpcOptions = append(grpcOptions, grpc.WithBlock(), grpc.WithAuthority("grpcuserapi.2dland.cn"), grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})), grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		// <!---- comment start ---->
		// check if accesstoken is expired, if expired, refresh it, this operation is not necessary for every request
		// it's just a demo, you should not do this in production environment
		// instead, you should refresh token when you get error code unauthenticated
		// or you can use a background goroutine/thread to refresh token periodically
		// thus it's not necessary, because the interceptor will refresh token automatically if token is expired
		//// ignoreAutoRefeshMethod := []string{pbPublicUser.PubUser_Login_FullMethodName, pbPublicUser.PubUser_Refresh_FullMethodName, pbPublicUser.PubUser_SendSmsVerifyCode_FullMethodName}
		////ignoreAutoRefesh := false
		////for _, m := range ignoreAutoRefeshMethod {
		////	if m == method {
		////		ignoreAutoRefesh = true
		////		break
		////	}
		////}
		////if !ignoreAutoRefesh && len(svc.accessToken) > 0 && accessTokenExpiredAt+120000 < time.Now().UnixMilli() && len(refreshToken) > 0 {
		////	// refresh token
		////	refreshResponse, err := pbPublicUser.NewPubUserClient(cc).Refresh(ctx, &pbPublicUser.Token{
		////		RefreshToken: refreshToken,
		////	})
		////	if err != nil {
		////		return err
		////	}
		////	if len(refreshResponse.AccessToken) > 0 {
		////		accessToken = refreshResponse.AccessToken
		////		accessTokenExpiredAt = refreshResponse.AccessTokenExpireTs
		////	}
		////}
		// <!---- comment end ---->
		// currentTimeStamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		ctxx := svc.signContext(method, ctx)
		err := invoker(ctxx, method, req, reply, cc, opts...) // invoking RPC method
		if err != nil {
			grpcStatus, ok := status.FromError(err)
			// if error is grpc error and error code is unauthenticated and error message contains "invalid accesstoken" and refresh token is not empty
			// then refresh access token and retry
			if ok && grpcStatus.Code() == codes.Unauthenticated && strings.Contains(grpcStatus.Err().Error(), "invalid accesstoken") && len(refreshToken) > 0 {
				// refresh token
				refreshResponse, err := pbPublicUser.NewPubUserClient(cc).Refresh(ctx, &pbPublicUser.Token{
					RefreshToken: refreshToken,
				})
				if err != nil {
					return err
				}
				if len(refreshResponse.AccessToken) > 0 {
					svc.accessToken = refreshResponse.AccessToken
					svc.accessTokenExpiredAt = refreshResponse.AccessTokenExpireTs
					svc.OnAccessTokenRefreshed(refreshResponse.AccessToken, refreshResponse.AccessTokenExpireTs, refreshResponse.RefreshToken, refreshResponse.RefreshTokenExpireTs)
				}
				// retry
				ctxx := svc.signContext(method, ctx)
				err = invoker(ctxx, method, req, reply, cc, opts...) // invoking RPC method
				if err != nil {
					return err
				} else {
					return nil
				}
			}
		}
		// post-processing
		return err
	}))
	grpcConnection, err := grpc.DialContext(dialContext, grpcServer, grpcOptions...)

	if err != nil {
		return nil, err
	}

	svc.grpcConnection = grpcConnection
	/*
		testCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		refreshResponse, err := pbPublicUser.NewPubUserClient(svc.grpcConnection).Refresh(testCtx, &pbPublicUser.Token{
			RefreshToken: refreshToken,
		})
		if err != nil {
			return nil, err
		}
		// if len(refreshResponse.AccessToken) > 0 {
		svc.OnAccessTokenRefreshed(refreshResponse.AccessToken, refreshResponse.AccessTokenExpireTs, refreshResponse.RefreshToken, refreshResponse.RefreshTokenExpireTs)
	*/

	return svc, err
}

func (s *AuthService) OnAccessTokenRefreshed(accessToken string, accessTokenExpiredAt int64, refreshToken string, refreshTokenExpiredAt int64) {
	// s.accessToken = accessToken
	// s.accessTokenExpiredAt = accessTokenExpiredAt
	s.configMutex.Lock()
	defer s.configMutex.Unlock()
	s.refreshToken = refreshToken
	viper.Set("refresh_token", refreshToken)
	// viper.Set("access_token", accessToken)
	// viper.Set("access_token_expired_at", accessTokenExpiredAt)
	viper.Set("refresh_token_expired_at", refreshTokenExpiredAt)
	err := viper.WriteConfig()
	if err != nil {
		err = viper.SafeWriteConfig()
	}
	if err != nil {
		println(err.Error())
	}

	if s.dopts.onTokenRefreshed != nil {
		s.dopts.onTokenRefreshed(accessToken, accessTokenExpiredAt, refreshToken, refreshTokenExpiredAt)
	}

}

func (s *AuthService) GetGrpcConnection() *grpc.ClientConn {
	return s.grpcConnection
}

func (s *AuthService) Close() {
	s.grpcConnection.Close()
}

func (s *AuthService) signContext(method string, ctx context.Context) context.Context {
	kvString := []string{}
	currentTimeStamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	bufferedString := bytes.NewBufferString(method)
	kvString = append(kvString, "timestamp", currentTimeStamp)
	bufferedString.WriteString(currentTimeStamp)
	kvString = append(kvString, "appid", s.appID)
	bufferedString.WriteString(s.appID)
	kvString = append(kvString, "appversion", s.appVersion)
	bufferedString.WriteString(s.appVersion)
	if len(s.accessToken) > 0 {
		authorization := "Bearer " + s.accessToken
		kvString = append(kvString, "authorization", authorization)
		bufferedString.WriteString(authorization)
	}
	bufferedString.WriteString(s.appSecret)
	sign := utils.GetMD5Hash(bufferedString.String())
	kvString = append(kvString, "sign", sign)
	return metadata.AppendToOutgoingContext(ctx, kvString...)
}
