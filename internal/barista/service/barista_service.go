package service

import (
	"coffee-shop/internal/barista/events"
	external_pub "coffee-shop/internal/barista/external/publisher"
	"coffee-shop/internal/barista/models"
	"coffee-shop/internal/barista/repository"
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

type baristaService struct {
	br repository.BaristaRepository
	cp external_pub.CounterEventPublisher
	v  validator.ValidatorInterface
	tx db.MongoTxInterface
}

func NewBaristaService(baristaRepo repository.BaristaRepository, v validator.ValidatorInterface, tx db.MongoTxInterface, cp external_pub.CounterEventPublisher) BaristaService {
	return &baristaService{br: baristaRepo, v: v, tx: tx, cp: cp}
}

func (b *baristaService) ProcessItems(ctx context.Context, req events.ItemsOrderedEvent) error {
	ctx, span := tracer.Start(ctx, "BaristaService.ProcessItems")
	defer span.End()
	err := b.v.Struct(req)
	if err != nil {
		return errors.BadRequest.Wrap(span, true, err, "Invalid request body")
	}
	fmt.Printf("Barista %#v\n", req)

	for _, item := range req.Items {
		id, err := b.br.CreateBaristaOrder(ctx, models.NewBaristaOrder(req.ProfileID, req.OrderID, item.OrderItemId, item.Type, item.Quantity))
		if err != nil {
			return err
		}
		fmt.Print(id.Hex())
		time.Sleep(2 * time.Second)
		err = b.br.UpdateBaristaOrder(ctx, bson.M{"_id": id}, bson.M{"item_status": constants.StatusFulfilled.String()})
		if err != nil {
			return err
		}

		eventBytes, err := json.Marshal(events.ItemOrderUpdated{OrderID: req.OrderID, ItemType: &item.Type, OrderItemID: item.OrderItemId, KitchenType: helpers.Ptr(int32(constants.Barista))})
		if err != nil {
			return errors.InternalError.Wrap(span, true, err, "")
		}
		err = b.cp.Publish(ctx, eventBytes, "text/plain")
		if err != nil {
			return errors.InternalError.Wrap(span, true, err, "Error while publishing message")
		}
	}

	if err != nil {
		return err
	}

	return err
}
