package service

import (
	"coffee-shop/internal/counter/events"
	external_grpc "coffee-shop/internal/counter/external/product"
	external_pub "coffee-shop/internal/counter/external/publisher"
	"coffee-shop/internal/counter/models"
	"coffee-shop/internal/counter/repository"
	"encoding/json"
	"log"
	"sync"

	"coffee-shop/pkg/constants"
	errors "coffee-shop/pkg/custom_errors"
	db "coffee-shop/pkg/db/mongo"
	"coffee-shop/pkg/validator"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("")

type counterService struct {
	cr    repository.CounterRepository
	p     external_grpc.ProductClient
	kp    external_pub.KitchenEventPublisher
	bp    external_pub.BaristaEventPublisher
	v     validator.ValidatorInterface
	tx    db.MongoTxInterface
	mutex sync.Mutex
}

func NewCounterService(counterRepo repository.CounterRepository, v validator.ValidatorInterface, tx db.MongoTxInterface, p external_grpc.ProductClient, kp external_pub.KitchenEventPublisher,
	bp external_pub.BaristaEventPublisher) CounterService {
	return &counterService{cr: counterRepo, v: v, tx: tx, p: p, kp: kp, bp: bp, mutex: sync.Mutex{}}
}

func (c *counterService) GetOrders(ctx context.Context, req models.GetOrdersReq) (*models.GetOrdersRes, error) {
	ctx, span := tracer.Start(ctx, "CounterService.PlaceOrder")
	defer span.End()
	err := c.v.Struct(req)
	if err != nil {
		return nil, errors.BadRequest.Wrap(span, true, err, "Invalid request body")
	}
	profileObjId, err := primitive.ObjectIDFromHex(req.ProfileID)
	if err != nil {
		return nil, errors.BadRequest.Wrap(span, true, err, "Invalid request body")
	}
	res, err := c.cr.GetOrders(ctx, bson.M{"profile_id": profileObjId})
	if err != nil {
		return nil, err
	}
	return &models.GetOrdersRes{
		Orders: res.GetOrderDetailsDTO(),
	}, nil
}

func (c *counterService) PlaceOrder(ctx context.Context, req models.PlaceOrderReq) error {
	ctx, span := tracer.Start(ctx, "CounterService.PlaceOrder")
	defer span.End()
	err := c.v.Struct(req)
	if err != nil {
		return errors.BadRequest.Wrap(span, true, err, "Invalid request body")
	}

	productMap, err := c.p.GetProductsByType(ctx, req.Items)
	if err != nil {
		return errors.InternalError.Wrap(span, true, err, "Error while calling products api")
	}

	if len(productMap) != len(req.Items) {
		return errors.BadRequest.New(span, true, "Invalid product type received!")
	}

	var totalPrice float32
	for _, p := range req.Items {
		totalPrice += productMap[*p.ItemType].Price * float32(p.Quantity)
	}

	orderId, err := c.cr.CreateOrder(ctx, models.NewOrder(
		req.ProfileID,
		totalPrice,
	))
	if err != nil {
		return err
	}

	kitchenEventMap := make(map[constants.KitchenType]events.ItemsOrderedEvent)

	for _, item := range req.Items {
		itemType := *item.ItemType
		itemId, err := c.cr.CreateOrderItem(ctx, models.NewOrderItem(orderId.Hex(), itemType, productMap[itemType].Price, item.Quantity))
		if err != nil {
			return err
		}
		itemKitchenType := productMap[itemType].KitchenType
		event, ok := kitchenEventMap[constants.KitchenType(itemKitchenType)]
		eventItem := events.Item{OrderItemId: itemId.Hex(), Type: itemType, Quantity: item.Quantity}
		if ok {
			event.Items = append(event.Items, eventItem)
			kitchenEventMap[constants.KitchenType(itemKitchenType)] = event
		} else {
			kitchenEventMap[constants.KitchenType(itemKitchenType)] = events.ItemsOrderedEvent{OrderID: orderId.Hex(), ProfileID: req.ProfileID, Items: []events.Item{eventItem}}
		}
	}

	baristEvent, ok := kitchenEventMap[constants.Barista]
	if ok {
		log.Print("B ", baristEvent)
		eventBytes, err := json.Marshal(baristEvent)
		if err != nil {
			return errors.InternalError.Wrap(span, true, err, "")
		}
		c.bp.Publish(ctx, eventBytes, "text/plain")
	}

	kitchenEvent, ok := kitchenEventMap[constants.Kitchen]
	if ok {
		log.Print("K", kitchenEvent)

		eventBytes, err := json.Marshal(kitchenEvent)
		if err != nil {
			return errors.InternalError.Wrap(span, true, err, "")
		}
		c.kp.Publish(ctx, eventBytes, "text/plain")
	}

	if err != nil {
		return err
	}

	return err
}

func (c *counterService) UpdateOrder(ctx context.Context, event events.ItemOrderUpdated) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	ctx, span := tracer.Start(ctx, "CounterService.UpdateOrder")
	defer span.End()
	err := c.v.Struct(event)
	if err != nil {
		return errors.BadRequest.Wrap(span, true, err, "Invalid request body")
	}

	orderObjectId, err := primitive.ObjectIDFromHex(event.OrderID)
	if err != nil {
		return errors.BadRequest.New(span, true, "Invalid order id received!")
	}
	allOrderItems, err := c.cr.GetOrderItems(ctx, bson.M{"order_id": orderObjectId})
	if err != nil {
		return err
	}

	orderItemObjectId, err := primitive.ObjectIDFromHex(event.OrderItemID)
	if err != nil {
		return errors.BadRequest.New(span, true, "Invalid order id received!")
	}
	orderItem, err := c.cr.GetOrderItem(ctx, bson.M{"_id": orderItemObjectId})
	if err != nil {
		return err
	}

	isOrderCompleted := true
	for _, o := range allOrderItems {
		if o.ItemStatus != constants.StatusFulfilled.String() && *o.Id != orderItemObjectId {
			isOrderCompleted = false
		}
	}

	if isOrderCompleted {
		err = c.cr.UpdateOrder(ctx, bson.M{"_id": orderItem.OrderId}, bson.M{"$set": bson.M{"order_status": constants.StatusFulfilled.String()}})
		if err != nil {
			return err
		}
	}

	err = c.cr.UpdateOrderItem(ctx, bson.M{"_id": orderItem.Id}, bson.M{"$set": bson.M{"item_status": constants.StatusFulfilled.String()}})
	if err != nil {
		return err
	}

	return nil
}
