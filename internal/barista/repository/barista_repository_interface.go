//go:generate mockgen -source account_repository_interface.go -destination ../mock/account_repository_mock.go -package mock
package repository

import (
	"coffee-shop/internal/barista/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaristaRepository interface {
	CreateBaristaOrder(ctx context.Context, order models.BaristaOrder) (*primitive.ObjectID, error)
	UpdateBaristaOrder(ctx context.Context, filter bson.M, update bson.M) error
}
