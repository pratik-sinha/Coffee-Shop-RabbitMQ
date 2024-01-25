//go:generate mockgen -source account_service_interface.go -destination ../mock/user_service_mock.go -package mock

package service

import (
	"coffee-shop/internal/barista/events"

	"golang.org/x/net/context"
)

type BaristaService interface {
	ProcessItems(context.Context, events.ItemsOrderedEvent) error
}
