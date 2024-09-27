package pubsub

import "context"

// PublisherSubscriberBroker interface is the simple message queue broker interface
// we will need for our application. Regardless of the implementation, all
// our application needs are the following three functions to operate.
type PubSubBroker interface {
	// Subscribe method allows a goroutine to subscribe to a topic.
	Subscribe(ctx context.Context, topicName string) <-chan []byte

	// Publish method allows a message to be published to a topic.
	Publish(ctx context.Context, topicName string, msg []byte) error

	// The Close method closes the agent and all channels in the subs map.
	Close()

	IsSubscriberConnectedToNetwork(ctx context.Context, topicName string) bool
}
