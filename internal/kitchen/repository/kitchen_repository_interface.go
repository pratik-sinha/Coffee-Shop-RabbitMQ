//go:generate mockgen -source account_repository_interface.go -destination ../mock/account_repository_mock.go -package mock
package repository

import (
	"coffee-shop/internal/kitchen/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KitchenRepository interface {
	CreateKitchenOrder(ctx context.Context, order models.KitchenOrder) (*primitive.ObjectID, error)
	UpdateKitchenOrder(ctx context.Context, filter bson.M, update bson.M) error
}
