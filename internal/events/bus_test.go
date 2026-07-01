package events

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Brook-sys/auxitalk/pkg/types"
)

func testEvent() types.Event {
	return types.Event{
		ID:        "event-1",
		Type:      "message.received",
		Source:    "test",
		CreatedAt: time.Now(),
	}
}

func TestPublishCallsTypedSubscriber(t *testing.T) {
	bus := New(Options{HandlerTimeout: time.Second})
	called := false

	_, err := bus.Subscribe("message.received", func(_ context.Context, event types.Event) error {
		called = true
		if event.Type != "message.received" {
			t.Fatalf("unexpected event type %q", event.Type)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	if err := bus.Publish(context.Background(), testEvent()); err != nil {
		t.Fatalf("publish: %v", err)
	}
	if !called {
		t.Fatal("expected handler to be called")
	}
}

func TestPublishCallsWildcardSubscriber(t *testing.T) {
	bus := New(Options{HandlerTimeout: time.Second})
	called := false

	_, err := bus.Subscribe("*", func(_ context.Context, _ types.Event) error {
		called = true
		return nil
	})
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	if err := bus.Publish(context.Background(), testEvent()); err != nil {
		t.Fatalf("publish: %v", err)
	}
	if !called {
		t.Fatal("expected wildcard handler to be called")
	}
}

func TestUnsubscribe(t *testing.T) {
	bus := New(Options{HandlerTimeout: time.Second})
	calls := 0

	sub, err := bus.Subscribe("message.received", func(_ context.Context, _ types.Event) error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	sub.Unsubscribe()
	if err := bus.Publish(context.Background(), testEvent()); err != nil {
		t.Fatalf("publish: %v", err)
	}
	if calls != 0 {
		t.Fatalf("expected 0 calls, got %d", calls)
	}
}

func TestPublishReturnsHandlerError(t *testing.T) {
	bus := New(Options{HandlerTimeout: time.Second})
	expected := errors.New("handler failed")

	_, err := bus.Subscribe("message.received", func(_ context.Context, _ types.Event) error {
		return expected
	})
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	if err := bus.Publish(context.Background(), testEvent()); !errors.Is(err, expected) {
		t.Fatalf("expected handler error, got %v", err)
	}
}

func TestPublishValidatesEvent(t *testing.T) {
	bus := New(Options{HandlerTimeout: time.Second})

	if err := bus.Publish(context.Background(), types.Event{}); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestHandlerTimeout(t *testing.T) {
	bus := New(Options{HandlerTimeout: 10 * time.Millisecond})

	_, err := bus.Subscribe("message.received", func(ctx context.Context, _ types.Event) error {
		<-ctx.Done()
		return ctx.Err()
	})
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	if err := bus.Publish(context.Background(), testEvent()); err == nil {
		t.Fatal("expected timeout error")
	}
}

func TestHistoryLimit(t *testing.T) {
	bus := New(Options{HandlerTimeout: time.Second, HistoryLimit: 2})

	event1 := testEvent()
	event1.ID = "event-1"
	event2 := testEvent()
	event2.ID = "event-2"
	event3 := testEvent()
	event3.ID = "event-3"

	for _, event := range []types.Event{event1, event2, event3} {
		if err := bus.Publish(context.Background(), event); err != nil {
			t.Fatalf("publish: %v", err)
		}
	}

	history := bus.History()
	if len(history) != 2 {
		t.Fatalf("expected 2 history events, got %d", len(history))
	}
	if history[0].ID != "event-2" || history[1].ID != "event-3" {
		t.Fatalf("unexpected history: %#v", history)
	}
}

func TestSubscribeRequiresHandler(t *testing.T) {
	bus := New(Options{})

	if _, err := bus.Subscribe("message.received", nil); err == nil {
		t.Fatal("expected missing handler error")
	}
}
