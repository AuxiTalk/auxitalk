package workflows

import (
	"errors"
	"testing"

	"github.com/Brook-sys/auxitalk/pkg/types"
)

func TestRegistryRegisterAndList(t *testing.T) {
	reg := NewRegistry()
	workflow := types.Workflow{
		ID:      "auto-reply",
		Enabled: true,
		Rules: []types.WorkflowRule{{
			ID:      "rule-1",
			Enabled: true,
			Trigger: types.WorkflowTrigger{EventType: "message.received"},
			Action:  types.WorkflowAction{Type: "message.reply.suggest", Risk: types.ActionRiskLow},
		}},
	}
	if err := reg.Register(workflow); err != nil {
		t.Fatalf("register: %v", err)
	}

	if len(reg.List()) != 1 {
		t.Fatalf("expected 1 workflow")
	}
	if _, ok := reg.Get("auto-reply"); !ok {
		t.Fatal("expected workflow to exist")
	}

	if err := reg.Register(workflow); !errors.Is(err, ErrWorkflowAlreadyExists) {
		t.Fatalf("expected already exists error, got %v", err)
	}
}

func TestRegistryEnableDisableAndEnabledRules(t *testing.T) {
	reg := NewRegistry()
	workflow := types.Workflow{
		ID:      "auto-reply",
		Enabled: false,
		Rules: []types.WorkflowRule{{
			ID:      "rule-1",
			Enabled: true,
			Trigger: types.WorkflowTrigger{EventType: "message.received"},
			Action:  types.WorkflowAction{Type: "message.reply.suggest", Risk: types.ActionRiskLow},
		}},
	}
	if err := reg.Register(workflow); err != nil {
		t.Fatalf("register: %v", err)
	}

	if len(reg.EnabledRules()) != 0 {
		t.Fatalf("expected 0 enabled rules when workflow disabled")
	}

	updated, err := reg.Enable("auto-reply")
	if err != nil || !updated.Enabled {
		t.Fatalf("enable: %v", err)
	}
	if len(reg.EnabledRules()) != 1 {
		t.Fatalf("expected 1 enabled rule")
	}

	_, err = reg.Disable("missing")
	if !errors.Is(err, ErrWorkflowNotFound) {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestWorkflowValidateRequiresRules(t *testing.T) {
	workflow := types.Workflow{ID: "empty"}
	if err := workflow.Validate(); err == nil {
		t.Fatal("expected validation error for empty workflow")
	}
}

func TestWorkflowDuplicateRuleID(t *testing.T) {
	workflow := types.Workflow{
		ID: "dup",
		Rules: []types.WorkflowRule{
			{ID: "rule-1", Enabled: true, Trigger: types.WorkflowTrigger{EventType: "message.received"}, Action: types.WorkflowAction{Type: "message.reply.suggest", Risk: types.ActionRiskLow}},
			{ID: "rule-1", Enabled: true, Trigger: types.WorkflowTrigger{EventType: "message.received"}, Action: types.WorkflowAction{Type: "message.reply.suggest", Risk: types.ActionRiskLow}},
		},
	}
	if err := workflow.Validate(); err == nil {
		t.Fatal("expected duplicate rule id error")
	}
}
