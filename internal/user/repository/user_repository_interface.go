//go:generate mockgen -source account_repository_interface.go -destination ../mock/account_repository_mock.go -package mock
package repository

import (
	"coffee-shop/internal/user/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) error
	GetUser(ctx context.Context, filter bson.M) (*models.User, error)
}
