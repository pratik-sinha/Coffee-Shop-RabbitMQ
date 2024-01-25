//go:generate mockgen -source account_service_interface.go -destination ../mock/user_service_mock.go -package mock

package service

import (
	"coffee-shop/internal/user/models"

	"golang.org/x/net/context"
)

type UserService interface {
	Register(context.Context, models.RegisterReq) error
	Login(context.Context, models.LoginReq) (*models.LoginRes, error)
	Logout(ctx context.Context, req models.LogoutReq) error
}
