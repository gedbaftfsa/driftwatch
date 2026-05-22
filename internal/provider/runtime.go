package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ServiceState represents the live runtime state of a deployed service.
type ServiceState struct {
	Name        string            `json:"name"`
	Image       string            `json:"image"`
	Replicas    int               `json:"replicas"`
	Env         map[string]string `json:"env"`
	LastUpdated time.Time         `json:"last_updated"`
}

// RuntimeProvider fetches live service state from a remote endpoint.
type RuntimeProvider struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewRuntimeProvider creates a RuntimeProvider with sensible defaults.
func NewRuntimeProvider(baseURL string) *RuntimeProvider {
	return &RuntimeProvider{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FetchAll retrieves the runtime state for all services from the provider endpoint.
func (r *RuntimeProvider) FetchAll(ctx context.Context) ([]ServiceState, error) {
	url := fmt.Sprintf("%s/services", r.BaseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching runtime state: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var states []ServiceState
	if err := json.NewDecoder(resp.Body).Decode(&states); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return states, nil
}

// FetchByName retrieves the runtime state for a single named service.
func (r *RuntimeProvider) FetchByName(ctx context.Context, name string) (*ServiceState, error) {
	url := fmt.Sprintf("%s/services/%s", r.BaseURL, name)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	resp, err := r.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching service %q: %w", name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var state ServiceState
	if err := json.NewDecoder(resp.Body).Decode(&state); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &state, nil
}
