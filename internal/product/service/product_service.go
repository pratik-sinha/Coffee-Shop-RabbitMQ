package service

import (
	"coffee-shop/internal/product/models"
	"coffee-shop/internal/product/repository"

	errors "coffee-shop/pkg/custom_errors"
	"coffee-shop/pkg/validator"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("")

type productService struct {
	pr repository.ProductRepository
	v  validator.ValidatorInterface
}

func NewProductService(productRepo repository.ProductRepository, v validator.ValidatorInterface) ProductService {
	return &productService{pr: productRepo, v: v}
}

func (p *productService) GetProducts(ctx context.Context) (*models.GetProductsRes, error) {
	ctx, span := tracer.Start(ctx, "ProductService.GetProducts")
	defer span.End()

	products, err := p.pr.GetProducts(ctx)
	if err != nil {
		return nil, err
	}

	return &models.GetProductsRes{
		Products: products.GetProductDto(),
	}, nil
}

func (p *productService) GetProductsByType(ctx context.Context, req models.GetProductsByTypeReq) (*models.GetProductsRes, error) {
	ctx, span := tracer.Start(ctx, "ProductService.GetProductsByType")
	defer span.End()
	err := p.v.Struct(req)
	if err != nil {
		return nil, errors.BadRequest.Wrap(span, true, err, "Invalid request body")
	}

	products, err := p.pr.GetProductsByType(ctx, bson.M{"type": bson.M{"$in": req.ProductTypes}})
	if err != nil {
		return nil, err
	}

	return &models.GetProductsRes{
		Products: products.GetProductDto(),
	}, nil

}
