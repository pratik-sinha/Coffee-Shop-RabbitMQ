package external

import (
	"context"
)

type BaristaEventPublisher interface {
	Configure()
	Publish(context.Context, []byte, string) error
}

type KitchenEventPublisher interface {
	Configure()
	Publish(context.Context, []byte, string) error
}
