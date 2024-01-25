//go:generate mockgen -source account_service_interface.go -destination ../mock/user_service_mock.go -package mock

package service

import (
	"coffee-shop/internal/product/models"

	"golang.org/x/net/context"
)

type ProductService interface {
	GetProducts(context.Context) (*models.GetProductsRes, error)
	GetProductsByType(ctx context.Context, req models.GetProductsByTypeReq) (*models.GetProductsRes, error)
}
