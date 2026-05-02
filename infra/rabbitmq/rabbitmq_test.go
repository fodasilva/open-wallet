package rabbitmq

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Skip test if no RabbitMQ running locally
func TestRabbitMQ(t *testing.T) {
	url := "amqp://guest:guest@localhost:5672/"
	client, err := NewClient(url)
	if err != nil {
		t.Skip("RabbitMQ not available at", url)
	}
	defer client.Close()

	queueName := "test-queue"
	_, err = client.DeclareQueue(queueName)
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	received := make(chan string)
	err = client.Consume(queueName, func(data []byte) error {
		received <- string(data)
		return nil
	})
	assert.NoError(t, err)

	payload := "hello rabbitmq"
	err = client.Publish(ctx, queueName, payload)
	assert.NoError(t, err)

	select {
	case msg := <-received:
		assert.Contains(t, msg, payload)
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for message")
	}
}
