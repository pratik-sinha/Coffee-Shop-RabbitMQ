package external

import (
	"coffee-shop/internal/counter/models"
	"coffee-shop/pkg/pb"
	"context"
)

type ProductClient interface {
	GetProductsByType(ctx context.Context, productTypes []models.ItemOrder) (map[int32]*pb.Product, error)
}
