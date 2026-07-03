package events

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Brook-sys/auxitalk/pkg/types"
)

type Handler func(context.Context, types.Event) error

type Options struct {
	HandlerTimeout time.Duration
	HistoryLimit   int
}

type Bus struct {
	mu             sync.RWMutex
	handlerTimeout time.Duration
	historyLimit   int
	subscribers    map[string]map[uint64]Handler
	wildcards      map[uint64]Handler
	history        []types.Event
	nextID         uint64
}

type Subscription struct {
	bus       *Bus
	eventType string
	id        uint64
	once      sync.Once
}

func New(options Options) *Bus {
	if options.HandlerTimeout <= 0 {
		options.HandlerTimeout = 5 * time.Second
	}

	return &Bus{
		handlerTimeout: options.HandlerTimeout,
		historyLimit:   options.HistoryLimit,
		subscribers:    make(map[string]map[uint64]Handler),
		wildcards:      make(map[uint64]Handler),
	}
}

func (b *Bus) Subscribe(eventType string, handler Handler) (*Subscription, error) {
	if handler == nil {
		return nil, errors.New("event handler is required")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.nextID++
	id := b.nextID

	if eventType == "" || eventType == "*" {
		b.wildcards[id] = handler
		return &Subscription{bus: b, eventType: "*", id: id}, nil
	}

	if b.subscribers[eventType] == nil {
		b.subscribers[eventType] = make(map[uint64]Handler)
	}
	b.subscribers[eventType][id] = handler

	return &Subscription{bus: b, eventType: eventType, id: id}, nil
}

func (s *Subscription) Unsubscribe() {
	if s == nil || s.bus == nil {
		return
	}

	s.once.Do(func() {
		s.bus.unsubscribe(s.eventType, s.id)
	})
}

func (b *Bus) Publish(ctx context.Context, event types.Event) error {
	if err := event.Validate(); err != nil {
		return err
	}

	handlers := b.snapshotHandlers(event.Type)
	b.record(event)

	var errs []error
	for _, handler := range handlers {
		if err := b.callHandler(ctx, handler, event); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (b *Bus) History() []types.Event {
	b.mu.RLock()
	defer b.mu.RUnlock()

	history := make([]types.Event, len(b.history))
	copy(history, b.history)
	return history
}

func (b *Bus) RestoreHistory(events []types.Event) {
	if b.historyLimit <= 0 {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.history = make([]types.Event, 0, len(events))
	for _, event := range events {
		b.history = append(b.history, event)
		if len(b.history) > b.historyLimit {
			b.history = b.history[len(b.history)-b.historyLimit:]
		}
	}
}

func (b *Bus) unsubscribe(eventType string, id uint64) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if eventType == "*" {
		delete(b.wildcards, id)
		return
	}

	delete(b.subscribers[eventType], id)
	if len(b.subscribers[eventType]) == 0 {
		delete(b.subscribers, eventType)
	}
}

func (b *Bus) snapshotHandlers(eventType string) []Handler {
	b.mu.RLock()
	defer b.mu.RUnlock()

	count := len(b.wildcards)
	if typed := b.subscribers[eventType]; typed != nil {
		count += len(typed)
	}

	handlers := make([]Handler, 0, count)
	for _, handler := range b.subscribers[eventType] {
		handlers = append(handlers, handler)
	}
	for _, handler := range b.wildcards {
		handlers = append(handlers, handler)
	}

	return handlers
}

func (b *Bus) record(event types.Event) {
	if b.historyLimit <= 0 {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.history = append(b.history, event)
	if len(b.history) > b.historyLimit {
		b.history = b.history[len(b.history)-b.historyLimit:]
	}
}

func (b *Bus) callHandler(ctx context.Context, handler Handler, event types.Event) error {
	handlerCtx, cancel := context.WithTimeout(ctx, b.handlerTimeout)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- handler(handlerCtx, event)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-handlerCtx.Done():
		return fmt.Errorf("event handler timeout for %s: %w", event.Type, handlerCtx.Err())
	case err := <-done:
		return err
	}
}
