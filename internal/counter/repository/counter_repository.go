package repository

import (
	"coffee-shop/internal/counter/models"
	"coffee-shop/pkg/custom_errors"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("")

type counterRepository struct {
	conn *mongo.Database
}

func NewCounterRepository(m *mongo.Database) CounterRepository {
	return &counterRepository{conn: m}
}

func (c *counterRepository) GetOrders(ctx context.Context, filter bson.M) (models.OrdersDetails, error) {
	ctx, span := tracer.Start(ctx, "CounterRepository.GetOrders")
	defer span.End()
	var res []models.OrderDetail
	pipelines := []bson.M{
		{"$match": filter},
		{"$lookup": bson.M{"from": "order_items", "localField": "_id", "foreignField": "order_id", "as": "order_items"}},
		{"$lookup": bson.M{"from": "products", "localField": "order_items.item_type", "foreignField": "type", "as": "product_details"}},
		{"$addFields": bson.M{"order_item_details": bson.M{"$map": bson.M{"input": "$order_items", "as": "o", "in": bson.M{"$mergeObjects": bson.A{"$$o", bson.M{"$arrayElemAt": bson.A{bson.M{"$filter": bson.M{"input": "$product_details", "as": "p", "cond": bson.M{"$eq": bson.A{"$$o.item_type", "$$p.type"}}}}, 0}}}}}}}},
		{"$project": bson.M{"_id": 0, "order_items": 0, "product_details": 0}},
	}

	cursor, err := c.conn.Collection("orders").Aggregate(context.TODO(), pipelines)
	if err != nil {
		return nil, custom_errors.InternalError.Wrapf(span, true, err, "Error while getting orders")
	}

	if err = cursor.All(ctx, &res); err != nil {
		return nil, custom_errors.InternalError.Wrapf(span, true, err, "Error while decoding orders")
	}
	return res, nil
}

func (c *counterRepository) CreateOrder(ctx context.Context, order models.Order) (*primitive.ObjectID, error) {
	ctx, span := tracer.Start(ctx, "CounterRepository.CreateOrder")
	defer span.End()
	res, err := c.conn.Collection("orders").InsertOne(ctx, order)
	if err != nil {
		err := custom_errors.InternalError.Wrapf(span, true, err, "Error while creating new order: %#v", order)
		return nil, err
	}
	id, _ := res.InsertedID.(primitive.ObjectID)
	return &id, nil
}

func (c *counterRepository) UpdateOrder(ctx context.Context, filter bson.M, update bson.M) error {
	fmt.Println("Ran")
	ctx, span := tracer.Start(ctx, "CounterRepository.UpdateOrder")
	defer span.End()
	res, err := c.conn.Collection("orders").UpdateOne(ctx, filter, update)
	if err != nil {
		err := custom_errors.InternalError.Wrapf(span, true, err, "Error while updating order filter: %#v update: %#v", filter, update)
		return err
	}
	if res.ModifiedCount == 1 {
		return nil
	} else {
		return custom_errors.InternalError.Wrapf(span, true, err, "No rows modified for updating order filter: %#v update: %#v", filter, update)
	}
}

func (c *counterRepository) UpdateOrderItem(ctx context.Context, filter bson.M, update bson.M) error {
	ctx, span := tracer.Start(ctx, "CounterRepository.UpdateOrderItem")
	defer span.End()
	res, err := c.conn.Collection("order_items").UpdateOne(ctx, filter, update)
	if err != nil {
		err := custom_errors.InternalError.Wrapf(span, true, err, "Error while updating order filter: %#v update: %#v", filter, update)
		return err
	}
	if res.ModifiedCount == 1 {
		return nil
	} else {

		return custom_errors.InternalError.Wrapf(span, true, err, "No rows modified for updating order filter: %#v update: %#v", filter, update)
	}
}

func (c *counterRepository) CreateOrderItem(ctx context.Context, item models.OrderItem) (*primitive.ObjectID, error) {
	ctx, span := tracer.Start(ctx, "CounterRepository.CreateOrderItem")
	defer span.End()
	res, err := c.conn.Collection("order_items").InsertOne(ctx, item)
	if err != nil {
		err := custom_errors.InternalError.Wrapf(span, true, err, "Error while creating new order item: %#v", item)
		return nil, err
	}
	id, _ := res.InsertedID.(primitive.ObjectID)
	return &id, nil
}

func (c *counterRepository) GetOrderItem(ctx context.Context, filter bson.M) (*models.OrderItem, error) {
	ctx, span := tracer.Start(ctx, "CounterRepository.GetOrderItem")
	defer span.End()
	models, err := c.GetOrderItems(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &models[0], nil
}

func (c *counterRepository) GetOrderItems(ctx context.Context, filter bson.M) ([]models.OrderItem, error) {
	ctx, span := tracer.Start(ctx, "CounterRepository.GetOrderItems")
	defer span.End()
	var res []models.OrderItem
	cursor, err := c.conn.Collection("order_items").Find(ctx, filter)
	if err != nil {
		err := custom_errors.InternalError.Wrapf(span, true, err, "Error while getting order item for filter: %#v", filter)
		return nil, err
	}
	err = cursor.All(ctx, &res)
	if err != nil {
		return nil, custom_errors.InternalError.Wrapf(span, true, err, "Error while decoding order itemsfor filter: %#v", filter)
	}
	if len(res) == 0 {
		return nil, custom_errors.InternalError.Newf(span, true, "No order  item found for filter:%#v", filter)
	}
	return res, nil
}
