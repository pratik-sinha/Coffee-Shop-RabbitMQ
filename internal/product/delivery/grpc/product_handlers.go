package grpc

import (
	"coffee-shop/internal/product/models"
	"coffee-shop/internal/product/service"
	errors "coffee-shop/pkg/custom_errors"
	"coffee-shop/pkg/pb"
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/protobuf/types/known/emptypb"
)

var tracer = otel.Tracer("")

type productHandler struct {
	pb.UnimplementedProductServiceServer
	ps service.ProductService
}

func NewProductHandler(productService service.ProductService) *productHandler {
	return &productHandler{ps: productService}
}

func (p *productHandler) GetProducts(ctx context.Context, req *emptypb.Empty) (*pb.GetProductsRes, error) {
	ctx, span := tracer.Start(ctx, "GetProducts.Controller.Grpc")
	defer span.End()
	span.SetAttributes(attribute.String("info", fmt.Sprintf("%#v", req)))
	res, err := p.ps.GetProducts(ctx)
	if err != nil {
		return nil, errors.HandleGrpcError(err)
	}
	return res.ConvertToProto(), nil
}

func (p *productHandler) GetProductsByType(ctx context.Context, req *pb.GetProductsByTypeReq) (*pb.GetProductsRes, error) {
	ctx, span := tracer.Start(ctx, "GetProductsByType.Controller.Grpc")
	defer span.End()
	var m models.GetProductsByTypeReq
	m.ConvertFromProto(req)
	span.SetAttributes(attribute.String("info", fmt.Sprintf("%#v", req)))
	res, err := p.ps.GetProductsByType(ctx, m)
	if err != nil {
		return nil, errors.HandleGrpcError(err)
	}
	return res.ConvertToProto(), nil
}
