package repository

import (
	"coffee-shop/internal/product/models"
	"coffee-shop/pkg/custom_errors"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("")

type productRepository struct {
	conn *mongo.Database
}

func NewProductRepository(m *mongo.Database) ProductRepository {
	return &productRepository{conn: m}
}

func (p *productRepository) GetProducts(ctx context.Context) (models.Products, error) {
	ctx, span := tracer.Start(ctx, "ProductRepository.GetProducts")
	defer span.End()

	var res []models.Product

	cursor, err := p.conn.Collection("products").Find(ctx, bson.M{})
	if err != nil {
		return nil, custom_errors.InternalError.Wrap(span, true, err, "Error while retrieving products")
	}

	err = cursor.All(ctx, &res)
	if err != nil {
		return nil, custom_errors.InternalError.Wrap(span, true, err, "Error while decoding products")
	}

	return res, nil
}

func (p *productRepository) GetProductsByType(ctx context.Context, filter bson.M) (models.Products, error) {
	ctx, span := tracer.Start(ctx, "ProductRepository.GetProductsByType")
	defer span.End()

	var res []models.Product

	cursor, err := p.conn.Collection("products").Find(ctx, filter)
	if err != nil {
		return nil, custom_errors.InternalError.Wrapf(span, true, err, "Error while retrieving products for filter: %#v", filter)
	}

	err = cursor.All(ctx, &res)
	if err != nil {
		return nil, custom_errors.InternalError.Wrapf(span, true, err, "Error while decoding products for filter: %#v", filter)
	}

	return res, nil
}
