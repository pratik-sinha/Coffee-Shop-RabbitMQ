package models

import (
	"coffee-shop/pkg/constants"
	"coffee-shop/pkg/helpers"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	Id          *primitive.ObjectID `bson:"_id,omitempty"`
	ProfileID   primitive.ObjectID  `bson:"profile_id"`
	OrderStatus string              `bson:"order_status"`
	TotalPrice  float32             `bson:"total_price"`
	CreatedAt   time.Time           `bson:"created_at"`
	UpdatedAt   time.Time           `bson:"updated_at"`
}

func NewOrder(profileID string, totalPrice float32) Order {
	objId, _ := primitive.ObjectIDFromHex(profileID)
	return Order{
		ProfileID:   objId,
		TotalPrice:  totalPrice,
		OrderStatus: constants.StatusInProcess.String(),
		CreatedAt:   helpers.GetUTCTimeStamp(),
		UpdatedAt:   helpers.GetUTCTimeStamp(),
	}
}

type OrderItem struct {
	Id         *primitive.ObjectID `bson:"_id,omitempty"`
	OrderId    primitive.ObjectID  `bson:"order_id"`
	ItemType   int32               `bson:"item_type"`
	Quantity   int32               `bson:"quantity"`
	Price      float32             `bson:"price"`
	ItemStatus string              `bson:"item_status"`
	CreatedAt  time.Time           `bson:"created_at"`
	UpdatedAt  time.Time           `bson:"updated_at"`
}

func NewOrderItem(orderId string, itemType int32, price float32, quantity int32) OrderItem {
	orderObjId, _ := primitive.ObjectIDFromHex(orderId)

	return OrderItem{
		OrderId:    orderObjId,
		ItemType:   itemType,
		Price:      price,
		Quantity:   quantity,
		ItemStatus: constants.StatusInProcess.String(),
		CreatedAt:  helpers.GetUTCTimeStamp(),
		UpdatedAt:  helpers.GetUTCTimeStamp(),
	}
}

//Agreggation Models

type OrderItemDetail struct {
	Name       string    `bson:"name"`
	Quantity   int32     `bson:"quantity"`
	Price      float32   `bson:"price"`
	ItemStatus string    `bson:"item_status"`
	CreatedAt  time.Time `bson:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at"`
}

type OrderItemsDetails []OrderItemDetail

func (o *OrderItemsDetails) GetOrderItemDTO() []OrderItemDTO {
	var res []OrderItemDTO
	for _, v := range *o {
		res = append(res, OrderItemDTO{
			Name:       v.Name,
			Quantity:   v.Quantity,
			Price:      v.Price,
			ItemStatus: v.ItemStatus,
		})
	}
	return res
}

type OrderDetail struct {
	OrderStatus string            `bson:"order_status"`
	TotalPrice  float32           `bson:"total_price"`
	CreatedAt   time.Time         `bson:"created_at"`
	UpdatedAt   time.Time         `bson:"updated_at"`
	OrderItems  OrderItemsDetails `bson:"order_item_details"`
}

type OrdersDetails []OrderDetail

func (o *OrdersDetails) GetOrderDetailsDTO() []OrderDTO {
	res := []OrderDTO{}
	for _, v := range *o {
		res = append(res, OrderDTO{
			OrderStatus: v.OrderStatus,
			OrderItems:  v.OrderItems.GetOrderItemDTO(),
			TotalPrice:  v.TotalPrice,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		})
	}
	return res
}
