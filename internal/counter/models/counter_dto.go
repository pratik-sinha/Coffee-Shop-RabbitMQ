package models

import "time"

type PlaceOrderReq struct {
	ProfileID string      `json:"profile_id" validate:"len=24"`
	Items     []ItemOrder `json:"items" validate:"min=1,dive"`
}

type GetOrdersReq struct {
	ProfileID string `json:"profile_id" validate:"len=24"`
}

type GetOrdersRes struct {
	Orders []OrderDTO `json:"orders"`
}

type OrderItemDTO struct {
	Name       string  `json:"name"`
	Quantity   int32   `json:"quantity"`
	Price      float32 `json:"price"`
	ItemStatus string  `json:"item_status"`
}

type OrderDTO struct {
	OrderStatus string         `json:"order_status"`
	OrderItems  []OrderItemDTO `json:"order_items"`
	TotalPrice  float32        `json:"total_price"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type ItemOrder struct {
	ItemType *int32 `json:"item_type" validate:"required"`
	Quantity int32  `json:"quantity" validate:"min=1"`
}
