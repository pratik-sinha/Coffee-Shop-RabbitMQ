package models

import (
	"coffee-shop/pkg/constants"
	"coffee-shop/pkg/helpers"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaristaOrder struct {
	Id          *primitive.ObjectID `bson:"_id,omitempty"`
	OrderID     string              `bson:"order_id"`
	OrderItemID string              `bson:"order_item_id"`
	ItemType    int32               `bson:"item_type"`
	Quantity    int32               `bson:"quantity"`
	ProfileID   string              `bson:"profile_id"`
	ItemStatus  string              `bson:"item_status"`
	CreatedAt   time.Time           `bson:"created_at"`
	UpdatedAt   time.Time           `bson:"updated_at"`
}

func NewBaristaOrder(profileID string, orderID string, orderItemID string, itemType int32, quantity int32) BaristaOrder {
	return BaristaOrder{
		ProfileID:   profileID,
		ItemType:    itemType,
		OrderID:     orderID,
		OrderItemID: orderItemID,
		Quantity:    quantity,
		ItemStatus:  constants.StatusInProcess.String(),
		CreatedAt:   helpers.GetUTCTimeStamp(),
		UpdatedAt:   helpers.GetUTCTimeStamp(),
	}
}
