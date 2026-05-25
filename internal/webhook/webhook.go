package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Payload is the JSON body sent to a webhook endpoint.
type Payload struct {
	Timestamp  time.Time        `json:"timestamp"`
	TotalDrift int              `json:"total_drift"`
	Summaries  []drift.Summary  `json:"summaries"`
}

// Sender dispatches drift summaries to a remote HTTP endpoint.
type Sender struct {
	client  *http.Client
	endpoint string
	secret  string
}

// NewSender creates a Sender that posts to endpoint.
// secret is optional; when non-empty it is sent as the
// X-Driftwatch-Secret header for basic authentication.
func NewSender(endpoint, secret string, timeout time.Duration) *Sender {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &Sender{
		client:   &http.Client{Timeout: timeout},
		endpoint: endpoint,
		secret:   secret,
	}
}

// Send serialises summaries and POSTs them to the configured endpoint.
// It returns an error if the HTTP response status is not 2xx.
func (s *Sender) Send(ctx context.Context, summaries []drift.Summary) error {
	payload := Payload{
		Timestamp:  time.Now().UTC(),
		TotalDrift: countDrifted(summaries),
		Summaries:  summaries,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if s.secret != "" {
		req.Header.Set("X-Driftwatch-Secret", s.secret)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d from %s", resp.StatusCode, s.endpoint)
	}
	return nil
}

func countDrifted(summaries []drift.Summary) int {
	n := 0
	for _, s := range summaries {
		if s.HasDrift {
			n++
		}
	}
	return n
}
