package distributedpubsub

import "context"

// PublishSubscribeBroker interface ....TODO
type PublishSubscribeBroker interface {
	// Subscribe method allows a goroutine to subscribe to a topic.
	Subscribe(ctx context.Context, topicName string) <-chan []byte

	// Publish method allows a message to be published to a topic.
	Publish(ctx context.Context, topicName string, msg []byte) error

	// The Close method closes the agent and all channels in the subs map.
	Close()
}
