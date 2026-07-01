package context

import (
	"strings"
	"testing"
	"time"

	"github.com/Brook-sys/auxitalk/pkg/types"
)

func TestBuilderBuild(t *testing.T) {
	builder := NewBuilder(10)
	session := types.Session{
		ID:      "s1",
		Channel: "mock",
		Messages: []types.Message{
			{ID: "m1", SessionID: "s1", AuthorID: "u1", Text: "hello", Direction: types.MessageDirectionInbound, Source: "test", CreatedAt: time.Now()},
			{ID: "m2", SessionID: "s1", AuthorID: "a1", Text: "hi", Direction: types.MessageDirectionOutbound, Source: "test", CreatedAt: time.Now()},
		},
	}

	ctx := builder.Build(session)
	if !strings.Contains(ctx, "User: hello") {
		t.Fatalf("expected user message, got %q", ctx)
	}
	if !strings.Contains(ctx, "Assistant: hi") {
		t.Fatalf("expected assistant message, got %q", ctx)
	}
}

func TestBuilderTruncates(t *testing.T) {
	builder := NewBuilder(2)
	session := types.Session{
		ID:      "s1",
		Channel: "mock",
		Messages: []types.Message{
			{ID: "m1", SessionID: "s1", AuthorID: "u1", Text: "old", Direction: types.MessageDirectionInbound, Source: "test", CreatedAt: time.Now()},
			{ID: "m2", SessionID: "s1", AuthorID: "u1", Text: "mid", Direction: types.MessageDirectionInbound, Source: "test", CreatedAt: time.Now()},
			{ID: "m3", SessionID: "s1", AuthorID: "u1", Text: "new", Direction: types.MessageDirectionInbound, Source: "test", CreatedAt: time.Now()},
		},
	}

	ctx := builder.Build(session)
	if strings.Contains(ctx, "old") {
		t.Fatalf("expected truncation of old message, got %q", ctx)
	}
	if !strings.Contains(ctx, "new") {
		t.Fatalf("expected newest message, got %q", ctx)
	}
}
