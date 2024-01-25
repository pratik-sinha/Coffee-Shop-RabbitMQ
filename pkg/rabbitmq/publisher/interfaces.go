package publisher

import "context"

type EventPublisher interface {
	Configure(...Option)
	Publish(context.Context, []byte, string) error
}
