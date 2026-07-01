package types

import (
	"errors"
	"strings"
	"time"
)

type ActionRisk string

const (
	ActionRiskLow    ActionRisk = "low"
	ActionRiskMedium ActionRisk = "medium"
	ActionRiskHigh   ActionRisk = "high"
)

type ActionStatus string

const (
	ActionStatusRequested ActionStatus = "requested"
	ActionStatusAllowed   ActionStatus = "allowed"
	ActionStatusConfirmed ActionStatus = "confirmed"
	ActionStatusDenied    ActionStatus = "denied"
	ActionStatusExecuted  ActionStatus = "executed"
	ActionStatusFailed    ActionStatus = "failed"
)

type ActionRequest struct {
	ID        string         `json:"id"`
	Type      string         `json:"type"`
	Risk      ActionRisk     `json:"risk"`
	Status    ActionStatus   `json:"status"`
	Source    string         `json:"source"`
	SessionID string         `json:"sessionId,omitempty"`
	Payload   map[string]any `json:"payload,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
}

func (a ActionRequest) Validate() error {
	if strings.TrimSpace(a.ID) == "" {
		return errors.New("action id is required")
	}
	if strings.TrimSpace(a.Type) == "" {
		return errors.New("action type is required")
	}
	if strings.TrimSpace(a.Source) == "" {
		return errors.New("action source is required")
	}
	switch a.Risk {
	case ActionRiskLow, ActionRiskMedium, ActionRiskHigh:
	default:
		return errors.New("action risk is invalid")
	}
	switch a.Status {
	case ActionStatusRequested, ActionStatusAllowed, ActionStatusConfirmed, ActionStatusDenied, ActionStatusExecuted, ActionStatusFailed:
	default:
		return errors.New("action status is invalid")
	}
	if a.CreatedAt.IsZero() {
		return errors.New("action createdAt is required")
	}
	return nil
}
