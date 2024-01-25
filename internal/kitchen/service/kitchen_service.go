package service

import (
	"coffee-shop/internal/kitchen/events"
	external_pub "coffee-shop/internal/kitchen/external/publisher"
	"coffee-shop/internal/kitchen/models"
	"coffee-shop/internal/kitchen/repository"
	"encoding/json"
	"fmt"
	"time"

	"coffee-shop/pkg/constants"
	errors "coffee-shop/pkg/custom_errors"
	db "coffee-shop/pkg/db/mongo"
	"coffee-shop/pkg/helpers"
	"coffee-shop/pkg/validator"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("")

type kitchenService struct {
	kr repository.KitchenRepository
	cp external_pub.CounterEventPublisher
	v  validator.ValidatorInterface
	tx db.MongoTxInterface
}

func NewKitchenService(kitchenRepo repository.KitchenRepository, v validator.ValidatorInterface, tx db.MongoTxInterface, cp external_pub.CounterEventPublisher) KitchenService {
	return &kitchenService{kr: kitchenRepo, v: v, tx: tx, cp: cp}
}

func (k *kitchenService) ProcessItems(ctx context.Context, req events.ItemsOrderedEvent) error {
	ctx, span := tracer.Start(ctx, "KitchenService.ProcessItems")
	defer span.End()
	err := k.v.Struct(req)
	if err != nil {
		return errors.BadRequest.Wrap(span, true, err, "Invalid request body")
	}
	fmt.Printf("Kitchen %#v\n", req)

	for _, item := range req.Items {
		id, err := k.kr.CreateKitchenOrder(ctx, models.NewKitchenOrder(req.ProfileID, req.OrderID, item.OrderItemId, item.Type, item.Quantity))
		if err != nil {
			return err
		}
		fmt.Print(id.Hex())
		time.Sleep(2 * time.Second)
		err = k.kr.UpdateKitchenOrder(ctx, bson.M{"_id": id}, bson.M{"item_status": constants.StatusFulfilled.String()})
		if err != nil {
			return err
		}

		eventBytes, err := json.Marshal(events.ItemOrderUpdated{OrderID: req.OrderID, ItemType: &item.Type, OrderItemID: item.OrderItemId, KitchenType: helpers.Ptr(int32(constants.Kitchen))})
		if err != nil {
			return errors.InternalError.Wrap(span, true, err, "")
		}
		err = k.cp.Publish(ctx, eventBytes, "text/plain")
		if err != nil {
			return errors.InternalError.Wrap(span, true, err, "Error while publishing message")
		}
	}

	if err != nil {
		return err
	}

	return err
}
