package events

type ItemsOrderedEvent struct {
	OrderID string
	Items   []Item
}

type Item struct {
	OrderItemId string
	Type        int32
	Quantity    int32
}

type ItemOrderUpdated struct {
	OrderID     string `validate:"len=24"`
	OrderItemID string `validate:"len=24"`
	KitchenType *int32 `validate:"required"`
	ItemType    *int32 `validate:"required"`
}
