package agent

import (
	"fmt"
	"sync"
)

// Registry manages agent registration and lookup.
type Registry struct {
	mu     sync.RWMutex
	agents map[string]Agent
}

// NewRegistry creates a new agent registry.
func NewRegistry() *Registry {
	return &Registry{
		agents: make(map[string]Agent),
	}
}

// Register adds an agent to the registry.
func (r *Registry) Register(a Agent) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.agents[a.Name()] = a
}

// Get retrieves an agent by name.
func (r *Registry) Get(name string) (Agent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.agents[name]
	if !ok {
		return nil, fmt.Errorf("agent not found: %s", name)
	}
	return a, nil
}

// All returns all registered agents.
func (r *Registry) All() map[string]Agent {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make(map[string]Agent, len(r.agents))
	for k, v := range r.agents {
		result[k] = v
	}
	return result
}
