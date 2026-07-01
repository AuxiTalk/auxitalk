package types

import (
	"errors"
	"strings"
)

type WorkflowTrigger struct {
	EventType string `json:"eventType"`
	Source    string `json:"source,omitempty"`
}

type WorkflowAction struct {
	Type    string         `json:"type"`
	Risk    ActionRisk     `json:"risk"`
	Payload map[string]any `json:"payload,omitempty"`
}

type Workflow struct {
	ID      string         `json:"id"`
	Name    string         `json:"name,omitempty"`
	Enabled bool           `json:"enabled"`
	Rules   []WorkflowRule `json:"rules"`
}

type WorkflowRule struct {
	ID      string          `json:"id"`
	Name    string          `json:"name,omitempty"`
	Enabled bool            `json:"enabled"`
	Trigger WorkflowTrigger `json:"trigger"`
	Action  WorkflowAction  `json:"action"`
}

func (w Workflow) Validate() error {
	if strings.TrimSpace(w.ID) == "" {
		return errors.New("workflow id is required")
	}
	if len(w.Rules) == 0 {
		return errors.New("workflow requires at least one rule")
	}
	seen := map[string]struct{}{}
	for _, rule := range w.Rules {
		if err := rule.Validate(); err != nil {
			return err
		}
		if _, exists := seen[rule.ID]; exists {
			return errors.New("workflow rule id must be unique")
		}
		seen[rule.ID] = struct{}{}
	}
	return nil
}

func (w Workflow) EnabledRules() []WorkflowRule {
	if !w.Enabled {
		return nil
	}
	rules := make([]WorkflowRule, 0, len(w.Rules))
	for _, rule := range w.Rules {
		if rule.Enabled {
			rules = append(rules, rule)
		}
	}
	return rules
}

func (r WorkflowRule) Validate() error {
	if strings.TrimSpace(r.ID) == "" {
		return errors.New("workflow rule id is required")
	}
	if strings.TrimSpace(r.Trigger.EventType) == "" {
		return errors.New("workflow trigger eventType is required")
	}
	if strings.TrimSpace(r.Action.Type) == "" {
		return errors.New("workflow action type is required")
	}
	switch r.Action.Risk {
	case ActionRiskLow, ActionRiskMedium, ActionRiskHigh:
	default:
		return errors.New("workflow action risk is invalid")
	}
	return nil
}

func (r WorkflowRule) Matches(event Event) bool {
	if !r.Enabled {
		return false
	}
	if r.Trigger.EventType != "*" && r.Trigger.EventType != event.Type {
		return false
	}
	if r.Trigger.Source != "" && r.Trigger.Source != event.Source {
		return false
	}
	return true
}
