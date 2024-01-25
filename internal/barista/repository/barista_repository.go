package repository

import (
	"coffee-shop/internal/barista/models"
	"coffee-shop/pkg/custom_errors"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("")

type baristaRepository struct {
	conn *mongo.Database
}

func NewBaristaRepository(m *mongo.Database) BaristaRepository {
	return &baristaRepository{conn: m}
}

func (b *baristaRepository) CreateBaristaOrder(ctx context.Context, order models.BaristaOrder) (*primitive.ObjectID, error) {
	ctx, span := tracer.Start(ctx, "BaristaRepository.CreateBaristaOrder")
	defer span.End()
	res, err := b.conn.Collection("barista_orders").InsertOne(ctx, order)
	if err != nil {
		err := custom_errors.InternalError.Wrapf(span, true, err, "Error while creating new barista order: %#v", order)
		return nil, err
	}
	id, _ := res.InsertedID.(primitive.ObjectID)
	return &id, nil
}

func (b *baristaRepository) UpdateBaristaOrder(ctx context.Context, filter bson.M, update bson.M) error {
	ctx, span := tracer.Start(ctx, "BaristaRepository.UpdateBaristaOrder")
	defer span.End()
	res, err := b.conn.Collection("barista_orders").UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		err := custom_errors.InternalError.Wrapf(span, true, err, "Error while updating barista order filter: %#v update: %#v", filter, update)
		return err
	}
	if res.ModifiedCount == 1 {
		return nil
	} else {
		return custom_errors.InternalError.Wrapf(span, true, err, "No rows modified for barista kitchen order filter: %#v update: %#v", filter, update)
	}
}
