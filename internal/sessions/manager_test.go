package sessions

import (
	"errors"
	"testing"
	"time"

	"github.com/Brook-sys/auxitalk/pkg/types"
)

func now() time.Time {
	return time.Now().UTC()
}

func TestManagerCreateGetUpdateDelete(t *testing.T) {
	m := NewManager()
	session := types.Session{
		ID:        "s1",
		Channel:   "mock",
		CreatedAt: now(),
		UpdatedAt: now(),
	}

	if err := m.Create(session); err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := m.Get("s1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.ID != "s1" {
		t.Fatalf("unexpected id: %s", got.ID)
	}

	session.Channel = "updated"
	if err := m.Update(session); err != nil {
		t.Fatalf("update: %v", err)
	}

	got, _ = m.Get("s1")
	if got.Channel != "updated" {
		t.Fatalf("channel not updated")
	}

	m.Delete("s1")
	if _, err := m.Get("s1"); !errors.Is(err, ErrSessionNotFound) {
		t.Fatalf("expected not found after delete, got %v", err)
	}
}

func TestManagerAppendMessage(t *testing.T) {
	m := NewManager()
	session := types.Session{
		ID:        "s1",
		Channel:   "mock",
		CreatedAt: now(),
		UpdatedAt: now(),
	}
	m.Create(session)

	msg := types.Message{
		ID:        "m1",
		SessionID: "s1",
		AuthorID:  "u1",
		Text:      "hello",
		Direction: types.MessageDirectionInbound,
		Source:    "test",
		CreatedAt: now(),
	}

	if err := m.AppendMessage("s1", msg); err != nil {
		t.Fatalf("append: %v", err)
	}

	got, _ := m.Get("s1")
	if len(got.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(got.Messages))
	}
}

func TestManagerRejectsDuplicate(t *testing.T) {
	m := NewManager()
	session := types.Session{ID: "s1", Channel: "mock", CreatedAt: now(), UpdatedAt: now()}
	m.Create(session)

	err := m.Create(session)
	if err == nil {
		t.Fatal("expected duplicate error")
	}
}
