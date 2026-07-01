package sessions

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Brook-sys/auxitalk/pkg/types"
)

var ErrSessionNotFound = errors.New("session not found")

type Manager struct {
	mu       sync.RWMutex
	sessions map[string]types.Session
}

func NewManager() *Manager {
	return &Manager{
		sessions: make(map[string]types.Session),
	}
}

func (m *Manager) Create(session types.Session) error {
	if err := session.Validate(); err != nil {
		return fmt.Errorf("invalid session: %w", err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.sessions[session.ID]; exists {
		return fmt.Errorf("session %s already exists", session.ID)
	}

	m.sessions[session.ID] = session
	return nil
}

func (m *Manager) Get(id string) (types.Session, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, ok := m.sessions[id]
	if !ok {
		return types.Session{}, ErrSessionNotFound
	}
	return session, nil
}

func (m *Manager) Update(session types.Session) error {
	if err := session.Validate(); err != nil {
		return fmt.Errorf("invalid session: %w", err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.sessions[session.ID]; !exists {
		return ErrSessionNotFound
	}

	m.sessions[session.ID] = session
	return nil
}

func (m *Manager) Delete(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.sessions, id)
}

func (m *Manager) List() []types.Session {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]types.Session, 0, len(m.sessions))
	for _, s := range m.sessions {
		result = append(result, s)
	}
	return result
}

func (m *Manager) AppendMessage(sessionID string, message types.Message) error {
	if err := message.Validate(); err != nil {
		return fmt.Errorf("invalid message: %w", err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	session, ok := m.sessions[sessionID]
	if !ok {
		return ErrSessionNotFound
	}

	session.Messages = append(session.Messages, message)
	session.UpdatedAt = time.Now()
	m.sessions[sessionID] = session
	return nil
}
