package services

import "context"

type BrokerService interface {
	Publish(ctx context.Context, queue string, body any) error
	Consume(queue string, handler func(data []byte) error) error
	Close() error
}
