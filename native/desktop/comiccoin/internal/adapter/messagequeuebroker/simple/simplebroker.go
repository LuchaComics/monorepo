// Package simple is a basic message queue implementation without any persistence, network functionality, nor anything more complex. This package takes advantage of the golang `goroutines` and provides a simple interface to use throughout your app.
package simple

import (
	"log"
	"log/slog"
	"sync"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/messagequeuebroker"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
)

type messageQueueBrokerImpl struct {
	// mu field is used to protect access to the subs and closed fields using a mutex.
	mu sync.Mutex

	// subs field is a map that holds a list of channels for each topic, allowing subscribers to receive messages published to that topic.
	subs map[string][]chan []byte

	// The quit field is a channel that is closed when the `message queue broker` is closed, allowing goroutines that are blocked on the channel to unblock and exit.
	quit chan struct{}

	// The closed field is a flag that indicates whether the `message queue broker` has been closed.
	closed bool
}

func NewMessageQueue(cfg *config.Config, logger *slog.Logger) messagequeuebroker.MessageQueueBroker {
	if cfg.DB.DataDir == "" {
		log.Fatal("cannot have empty dir")
	}

	return &messageQueueBrokerImpl{
		subs: make(map[string][]chan []byte),
		quit: make(chan struct{}),
	}
}

func (impl *messageQueueBrokerImpl) Subscribe(topic string) <-chan []byte {
	impl.mu.Lock()
	defer impl.mu.Unlock()

	if impl.closed {
		return nil
	}

	ch := make(chan []byte)
	impl.subs[topic] = append(impl.subs[topic], ch)
	return ch
}

func (impl *messageQueueBrokerImpl) Publish(topic string, msg []byte) {
	impl.mu.Lock()
	defer impl.mu.Unlock()

	if impl.closed {
		return
	}

	for _, ch := range impl.subs[topic] {
		ch <- msg
	}
}

func (impl *messageQueueBrokerImpl) Close() {
	impl.mu.Lock()
	defer impl.mu.Unlock()

	if impl.closed {
		return
	}

	impl.closed = true
	close(impl.quit)

	for _, ch := range impl.subs {
		for _, sub := range ch {
			close(sub)
		}
	}
}
