package distributedmessagequeue

import "context"

// DistributedMessageQueueBroker interface ....TODO
type DistributedMessageQueueBroker interface {
	// Subscribe method allows a goroutine to subscribe to a topic.
	Subscribe(ctx context.Context, topicName string) []byte

	// Publish method allows a message to be published to a topic.
	Publish(ctx context.Context, topicName string, msg []byte)

	// The Close method closes the agent and all channels in the subs map.
	Close()
}
