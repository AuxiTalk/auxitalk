package types

import (
	"errors"
	"strings"
	"time"
)

type SuggestionAction string

const (
	SuggestionActionRespond SuggestionAction = "respond"
	SuggestionActionWait    SuggestionAction = "wait"
	SuggestionActionAskUser SuggestionAction = "ask_user"
	SuggestionActionIgnore  SuggestionAction = "ignore"
)

type Suggestion struct {
	ID         string           `json:"id"`
	SessionID  string           `json:"sessionId"`
	Text       string           `json:"text,omitempty"`
	Confidence float64          `json:"confidence"`
	Reason     string           `json:"reason,omitempty"`
	Tone       string           `json:"tone,omitempty"`
	Action     SuggestionAction `json:"action"`
	Metadata   map[string]any   `json:"metadata,omitempty"`
	CreatedAt  time.Time        `json:"createdAt"`
}

func (s Suggestion) Validate() error {
	if strings.TrimSpace(s.ID) == "" {
		return errors.New("suggestion id is required")
	}
	if strings.TrimSpace(s.SessionID) == "" {
		return errors.New("suggestion sessionId is required")
	}
	if s.Confidence < 0 || s.Confidence > 1 {
		return errors.New("suggestion confidence must be between 0 and 1")
	}
	switch s.Action {
	case SuggestionActionRespond, SuggestionActionWait, SuggestionActionAskUser, SuggestionActionIgnore:
	default:
		return errors.New("suggestion action is invalid")
	}
	if s.Action == SuggestionActionRespond && strings.TrimSpace(s.Text) == "" {
		return errors.New("suggestion text is required for respond action")
	}
	if s.CreatedAt.IsZero() {
		return errors.New("suggestion createdAt is required")
	}
	return nil
}
