package types

import (
	"testing"
	"time"
)

func TestEventValidate(t *testing.T) {
	event := Event{
		ID:        "event-1",
		Type:      "message.received",
		Source:    "mock-input",
		CreatedAt: time.Now(),
	}

	if err := event.Validate(); err != nil {
		t.Fatalf("expected valid event: %v", err)
	}

	event.ID = ""
	if err := event.Validate(); err == nil {
		t.Fatal("expected missing id error")
	}
}

func TestMessageValidate(t *testing.T) {
	message := Message{
		ID:        "message-1",
		SessionID: "session-1",
		AuthorID:  "user-1",
		Text:      "hello",
		Direction: MessageDirectionInbound,
		Source:    "mock-input",
		CreatedAt: time.Now(),
	}

	if err := message.Validate(); err != nil {
		t.Fatalf("expected valid message: %v", err)
	}

	message.Direction = "sideways"
	if err := message.Validate(); err == nil {
		t.Fatal("expected invalid direction error")
	}
}

func TestSessionValidate(t *testing.T) {
	now := time.Now()
	session := Session{
		ID:        "session-1",
		Channel:   "mock",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := session.Validate(); err != nil {
		t.Fatalf("expected valid session: %v", err)
	}

	session.Channel = ""
	if err := session.Validate(); err == nil {
		t.Fatal("expected missing channel error")
	}
}

func TestSuggestionValidate(t *testing.T) {
	suggestion := Suggestion{
		ID:         "suggestion-1",
		SessionID:  "session-1",
		Text:       "hello back",
		Confidence: 0.8,
		Action:     SuggestionActionRespond,
		CreatedAt:  time.Now(),
	}

	if err := suggestion.Validate(); err != nil {
		t.Fatalf("expected valid suggestion: %v", err)
	}

	suggestion.Confidence = 1.5
	if err := suggestion.Validate(); err == nil {
		t.Fatal("expected invalid confidence error")
	}
}

func TestActionRequestValidate(t *testing.T) {
	action := ActionRequest{
		ID:        "action-1",
		Type:      "message.send",
		Risk:      ActionRiskHigh,
		Status:    ActionStatusRequested,
		Source:    "console-output",
		CreatedAt: time.Now(),
	}

	if err := action.Validate(); err != nil {
		t.Fatalf("expected valid action: %v", err)
	}

	action.Risk = "unknown"
	if err := action.Validate(); err == nil {
		t.Fatal("expected invalid risk error")
	}
}

func TestPluginManifestValidate(t *testing.T) {
	manifest := PluginManifest{
		ID:      "mock-input",
		Name:    "Mock Input",
		Version: "0.1.0",
		Runtime: "node",
		Entry:   "index.js",
		Kind:    PluginKindInput,
		Permissions: []string{
			"event.emit",
		},
		Capabilities: []Capability{
			{Name: "conversation.observe"},
		},
	}

	if err := manifest.Validate(); err != nil {
		t.Fatalf("expected valid manifest: %v", err)
	}

	manifest.Capabilities[0].Name = ""
	if err := manifest.Validate(); err == nil {
		t.Fatal("expected invalid capability error")
	}
}
