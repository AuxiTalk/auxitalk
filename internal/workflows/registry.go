package workflows

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Brook-sys/auxitalk/pkg/types"
)

var (
	ErrWorkflowNotFound      = errors.New("workflow not found")
	ErrWorkflowAlreadyExists = errors.New("workflow already exists")
)

type Registry struct {
	mu        sync.RWMutex
	workflows map[string]types.Workflow
}

func NewRegistry() *Registry {
	return &Registry{workflows: map[string]types.Workflow{}}
}

func (r *Registry) Register(workflow types.Workflow) error {
	if err := workflow.Validate(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.workflows[workflow.ID]; exists {
		return fmt.Errorf("%w: %s", ErrWorkflowAlreadyExists, workflow.ID)
	}
	r.workflows[workflow.ID] = workflow
	return nil
}

func (r *Registry) List() []types.Workflow {
	r.mu.RLock()
	defer r.mu.RUnlock()
	workflows := make([]types.Workflow, 0, len(r.workflows))
	for _, workflow := range r.workflows {
		workflows = append(workflows, workflow)
	}
	return workflows
}

func (r *Registry) Get(id string) (types.Workflow, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	workflow, ok := r.workflows[id]
	return workflow, ok
}

func (r *Registry) Enable(id string) (types.Workflow, error) {
	return r.setEnabled(id, true)
}

func (r *Registry) Disable(id string) (types.Workflow, error) {
	return r.setEnabled(id, false)
}

func (r *Registry) EnabledRules() []types.WorkflowRule {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rules := []types.WorkflowRule{}
	for _, workflow := range r.workflows {
		rules = append(rules, workflow.EnabledRules()...)
	}
	return rules
}

func (r *Registry) setEnabled(id string, enabled bool) (types.Workflow, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	workflow, ok := r.workflows[id]
	if !ok {
		return types.Workflow{}, fmt.Errorf("%w: %s", ErrWorkflowNotFound, id)
	}
	workflow.Enabled = enabled
	r.workflows[id] = workflow
	return workflow, nil
}
