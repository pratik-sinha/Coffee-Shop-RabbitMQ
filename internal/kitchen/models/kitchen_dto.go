package models

type PlaceOrderReq struct {
	ProfileID string      `json:"profile_id" validate:"len=24"`
	Items     []ItemOrder `json:"items" validate:"min=1,dive"`
}

type ItemOrder struct {
	ItemType *int32 `json:"item_type" validate:"required"`
	Quantity int32  `json:"quantity" validate:"min=1"`
}
