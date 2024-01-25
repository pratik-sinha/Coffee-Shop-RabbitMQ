package grpc

import (
	"coffee-shop/pkg/pb"
	"context"
)

type ProductController interface {
	GetProducts(ctx context.Context) (*pb.GetProductsRes, error)
	GetProductsByType(ctx context.Context, req pb.GetProductsByTypeReq) (*pb.GetProductsRes, error)
}
