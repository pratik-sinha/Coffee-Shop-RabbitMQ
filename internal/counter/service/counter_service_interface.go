//go:generate mockgen -source account_service_interface.go -destination ../mock/user_service_mock.go -package mock

package service

import (
	"coffee-shop/internal/counter/events"
	"coffee-shop/internal/counter/models"

	"golang.org/x/net/context"
)

type CounterService interface {
	GetOrders(ctx context.Context, req models.GetOrdersReq) (*models.GetOrdersRes, error)
	PlaceOrder(context.Context, models.PlaceOrderReq) error
	UpdateOrder(context.Context, events.ItemOrderUpdated) error
}
