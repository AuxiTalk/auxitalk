package types

import (
	"errors"
	"strings"
	"time"
)

type Event struct {
	ID        string         `json:"id"`
	Type      string         `json:"type"`
	Source    string         `json:"source"`
	SessionID string         `json:"sessionId,omitempty"`
	Payload   map[string]any `json:"payload,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
}

func (e Event) Validate() error {
	if strings.TrimSpace(e.ID) == "" {
		return errors.New("event id is required")
	}
	if strings.TrimSpace(e.Type) == "" {
		return errors.New("event type is required")
	}
	if strings.TrimSpace(e.Source) == "" {
		return errors.New("event source is required")
	}
	if e.CreatedAt.IsZero() {
		return errors.New("event createdAt is required")
	}
	return nil
}
