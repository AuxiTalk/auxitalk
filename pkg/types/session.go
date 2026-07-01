package types

import (
	"errors"
	"strings"
	"time"
)

type MessageDirection string

const (
	MessageDirectionInbound  MessageDirection = "inbound"
	MessageDirectionOutbound MessageDirection = "outbound"
)

type Participant struct {
	ID          string         `json:"id"`
	DisplayName string         `json:"displayName,omitempty"`
	Role        string         `json:"role,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

type Message struct {
	ID        string           `json:"id"`
	SessionID string           `json:"sessionId"`
	AuthorID  string           `json:"authorId"`
	Text      string           `json:"text"`
	Direction MessageDirection `json:"direction"`
	Source    string           `json:"source"`
	Metadata  map[string]any   `json:"metadata,omitempty"`
	CreatedAt time.Time        `json:"createdAt"`
}

func (m Message) Validate() error {
	if strings.TrimSpace(m.ID) == "" {
		return errors.New("message id is required")
	}
	if strings.TrimSpace(m.SessionID) == "" {
		return errors.New("message sessionId is required")
	}
	if strings.TrimSpace(m.AuthorID) == "" {
		return errors.New("message authorId is required")
	}
	if strings.TrimSpace(m.Source) == "" {
		return errors.New("message source is required")
	}
	if m.Direction != MessageDirectionInbound && m.Direction != MessageDirectionOutbound {
		return errors.New("message direction is invalid")
	}
	if m.CreatedAt.IsZero() {
		return errors.New("message createdAt is required")
	}
	return nil
}

type Session struct {
	ID           string         `json:"id"`
	Channel      string         `json:"channel"`
	Participants []Participant  `json:"participants,omitempty"`
	Messages     []Message      `json:"messages,omitempty"`
	State        string         `json:"state,omitempty"`
	Metadata     map[string]any `json:"metadata,omitempty"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
}

func (s Session) Validate() error {
	if strings.TrimSpace(s.ID) == "" {
		return errors.New("session id is required")
	}
	if strings.TrimSpace(s.Channel) == "" {
		return errors.New("session channel is required")
	}
	if s.CreatedAt.IsZero() {
		return errors.New("session createdAt is required")
	}
	if s.UpdatedAt.IsZero() {
		return errors.New("session updatedAt is required")
	}
	return nil
}
