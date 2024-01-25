package repository

import (
	"coffee-shop/internal/kitchen/models"
	"coffee-shop/pkg/custom_errors"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("")

type kitchenRepository struct {
	conn *mongo.Database
}

func NewKitchenRepository(m *mongo.Database) KitchenRepository {
	return &kitchenRepository{conn: m}
}

func (k *kitchenRepository) CreateKitchenOrder(ctx context.Context, order models.KitchenOrder) (*primitive.ObjectID, error) {
	ctx, span := tracer.Start(ctx, "CreateKitchenOrder.CreateOrder")
	defer span.End()
	res, err := k.conn.Collection("kitchen_orders").InsertOne(ctx, order)
	if err != nil {
		err := custom_errors.InternalError.Wrapf(span, true, err, "Error while creating new order: %#v", order)
		return nil, err
	}
	id, _ := res.InsertedID.(primitive.ObjectID)
	return &id, nil
}

func (k *kitchenRepository) UpdateKitchenOrder(ctx context.Context, filter bson.M, update bson.M) error {
	ctx, span := tracer.Start(ctx, "CreateKitchenOrder.UpdateKitchenOrder")
	defer span.End()
	res, err := k.conn.Collection("kitchen_orders").UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		err := custom_errors.InternalError.Wrapf(span, true, err, "Error while updating kitchen order filter: %#v update: %#v", filter, update)
		return err
	}
	if res.ModifiedCount == 1 {
		return nil
	} else {
		return custom_errors.InternalError.Wrapf(span, true, err, "No rows modified for updating kitchen order filter: %#v update: %#v", filter, update)
	}
}
