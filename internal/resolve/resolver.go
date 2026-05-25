// Package resolve maps declared service names to live runtime states,
// producing a joined view used by the drift detector.
package resolve

import (
	"context"
	"fmt"

	"github.com/yourorg/driftwatch/internal/config"
	"github.com/yourorg/driftwatch/internal/provider"
)

// ServiceState holds both the declared config entry and the live runtime state
// for a single service. Live may be nil when the service is not found at runtime.
type ServiceState struct {
	Name     string
	Declared config.Service
	Live     *provider.ServiceState
}

// Resolver joins declared config services with live runtime states.
type Resolver struct {
	runtime *provider.RuntimeProvider
}

// New creates a Resolver backed by the given RuntimeProvider.
func New(rt *provider.RuntimeProvider) *Resolver {
	return &Resolver{runtime: rt}
}

// Resolve fetches all live states and joins them with the declared services.
// Services present in config but absent at runtime have a nil Live field.
// Services present at runtime but absent in config are silently ignored.
func (r *Resolver) Resolve(ctx context.Context, declared []config.Service) ([]ServiceState, error) {
	liveList, err := r.runtime.FetchAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("resolve: fetch runtime states: %w", err)
	}

	liveMap := make(map[string]*provider.ServiceState, len(liveList))
	for i := range liveList {
		liveMap[liveList[i].Name] = &liveList[i]
	}

	states := make([]ServiceState, 0, len(declared))
	for _, svc := range declared {
		ss := ServiceState{
			Name:     svc.Name,
			Declared: svc,
			Live:     liveMap[svc.Name],
		}
		states = append(states, ss)
	}
	return states, nil
}
