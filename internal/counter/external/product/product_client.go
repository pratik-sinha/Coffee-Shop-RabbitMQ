package external

import (
	"coffee-shop/config"
	"coffee-shop/internal/counter/models"
	"coffee-shop/pkg/pb"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type productGRPCClient struct {
	conn *grpc.ClientConn
}

func NewGRPCProductClient(cfg *config.Config) (ProductClient, error) {
	conn, err := grpc.Dial(cfg.Product_service.GrpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &productGRPCClient{
		conn: conn,
	}, nil
}

func (p *productGRPCClient) GetProductsByType(ctx context.Context, products []models.ItemOrder) (map[int32]*pb.Product, error) {
	c := pb.NewProductServiceClient(p.conn)

	productTypes := []int32{}

	for _, v := range products {
		productTypes = append(productTypes, *v.ItemType)
	}

	res, err := c.GetProductsByType(ctx, &pb.GetProductsByTypeReq{ProductTypes: productTypes})
	if err != nil {
		return nil, err
	}

	productMap := make(map[int32]*pb.Product, len(res.Products))

	for _, p := range res.Products {
		productMap[p.Type] = p
	}

	return productMap, nil
}
