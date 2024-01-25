package external

import (
	"context"
)

type CounterEventPublisher interface {
	Configure()
	Publish(context.Context, []byte, string) error
}
