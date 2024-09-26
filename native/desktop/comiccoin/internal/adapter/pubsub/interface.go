package pubsub

// PubSubBroker interface is the simple message queue broker interface
// we will need for our application. Regardless of the implementation, all
// our application needs are the following three functions to operate.
type PublisherSubscriberBroker interface {
	// Subscribe method allows a goroutine to subscribe to a topic.
	Subscribe(topic string) <-chan []byte

	// Publish method allows a message to be published to a topic.
	Publish(topic string, msg []byte)

	// The Close method closes the agent and all channels in the subs map.
	Close()
}
