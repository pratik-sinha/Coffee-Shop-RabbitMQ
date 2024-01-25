//go:generate mockgen -source account_repository_interface.go -destination ../mock/account_repository_mock.go -package mock
package repository

import (
	"coffee-shop/internal/counter/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CounterRepository interface {
	GetOrders(ctx context.Context, filter bson.M) (models.OrdersDetails, error)
	CreateOrder(ctx context.Context, order models.Order) (*primitive.ObjectID, error)
	UpdateOrder(ctx context.Context, filter bson.M, update bson.M) error
	GetOrderItem(ctx context.Context, filter bson.M) (*models.OrderItem, error)
	GetOrderItems(ctx context.Context, filter bson.M) ([]models.OrderItem, error)
	CreateOrderItem(ctx context.Context, item models.OrderItem) (*primitive.ObjectID, error)
	UpdateOrderItem(ctx context.Context, filter bson.M, update bson.M) error
}
