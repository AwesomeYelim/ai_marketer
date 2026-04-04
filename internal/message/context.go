package message

import "sync"

// WorkflowContext holds shared state across agent executions in a pipeline.
// It is concurrency-safe via a read-write mutex.
type WorkflowContext struct {
	mu sync.RWMutex

	// UserRequest is the original user input.
	UserRequest string `json:"user_request"`

	// Results stores each agent's output keyed by agent name.
	Results map[string]*AgentResponse `json:"results"`

	// SharedData stores arbitrary key-value data shared across agents.
	SharedData map[string]interface{} `json:"shared_data"`

	// RetryCount tracks the number of retries for the trust loop.
	RetryCount int `json:"retry_count"`

	// MaxRetries is the maximum number of trust-loop retries allowed.
	MaxRetries int `json:"max_retries"`
}

// NewWorkflowContext creates a new WorkflowContext.
func NewWorkflowContext(userRequest string) *WorkflowContext {
	return &WorkflowContext{
		UserRequest: userRequest,
		Results:     make(map[string]*AgentResponse),
		SharedData:  make(map[string]interface{}),
		MaxRetries:  2,
	}
}

// SetResult stores an agent's response.
func (wc *WorkflowContext) SetResult(agentName string, resp *AgentResponse) {
	wc.mu.Lock()
	defer wc.mu.Unlock()
	wc.Results[agentName] = resp
}

// GetResult retrieves an agent's response.
func (wc *WorkflowContext) GetResult(agentName string) (*AgentResponse, bool) {
	wc.mu.RLock()
	defer wc.mu.RUnlock()
	r, ok := wc.Results[agentName]
	return r, ok
}

// Set stores a value in shared data.
func (wc *WorkflowContext) Set(key string, value interface{}) {
	wc.mu.Lock()
	defer wc.mu.Unlock()
	wc.SharedData[key] = value
}

// Get retrieves a value from shared data.
func (wc *WorkflowContext) Get(key string) (interface{}, bool) {
	wc.mu.RLock()
	defer wc.mu.RUnlock()
	v, ok := wc.SharedData[key]
	return v, ok
}

// IncrementRetry increments the retry counter and returns true if under the limit.
func (wc *WorkflowContext) IncrementRetry() bool {
	wc.mu.Lock()
	defer wc.mu.Unlock()
	wc.RetryCount++
	return wc.RetryCount <= wc.MaxRetries
}
