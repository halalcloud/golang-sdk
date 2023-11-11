package user

import "github.com/halalcloud/golang-sdk/auth"

type UserService struct {
	authService *auth.AuthService
}

func NewUserService(authService *auth.AuthService) *UserService {
	return &UserService{
		authService: authService,
	}
}
