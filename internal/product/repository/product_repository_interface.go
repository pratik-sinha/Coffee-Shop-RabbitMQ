//go:generate mockgen -source account_repository_interface.go -destination ../mock/account_repository_mock.go -package mock
package repository

import (
	"coffee-shop/internal/product/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type ProductRepository interface {
	GetProducts(ctx context.Context) (models.Products, error)
	GetProductsByType(ctx context.Context, filter bson.M) (models.Products, error)
}
