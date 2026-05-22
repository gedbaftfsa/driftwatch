package drift

import (
	"fmt"
	"strings"
)

// DriftResult holds the result of a drift check for a single service.
type DriftResult struct {
	ServiceName string
	Drifted     bool
	Diffs       []string
}

// ServiceState represents the declared state of a service from config.
type ServiceState struct {
	Name    string
	Image   string
	Replicas int
	Env     map[string]string
}

// DeployedState represents the actual running state of a service.
type DeployedState struct {
	Name    string
	Image   string
	Replicas int
	Env     map[string]string
}

// Detect compares declared service states against deployed states and returns drift results.
func Detect(declared []ServiceState, deployed map[string]DeployedState) ([]DriftResult, error) {
	if declared == nil {
		return nil, fmt.Errorf("declared service states must not be nil")
	}

	results := make([]DriftResult, 0, len(declared))

	for _, svc := range declared {
		result := DriftResult{ServiceName: svc.Name}

		actual, found := deployed[svc.Name]
		if !found {
			result.Drifted = true
			result.Diffs = append(result.Diffs, fmt.Sprintf("service %q not found in deployed state", svc.Name))
			results = append(results, result)
			continue
		}

		if svc.Image != actual.Image {
			result.Diffs = append(result.Diffs, fmt.Sprintf("image: declared=%q actual=%q", svc.Image, actual.Image))
		}

		if svc.Replicas != actual.Replicas {
			result.Diffs = append(result.Diffs, fmt.Sprintf("replicas: declared=%d actual=%d", svc.Replicas, actual.Replicas))
		}

		for k, declaredVal := range svc.Env {
			actualVal, ok := actual.Env[k]
			if !ok {
				result.Diffs = append(result.Diffs, fmt.Sprintf("env %q: declared=%q actual=<missing>", k, declaredVal))
			} else if declaredVal != actualVal {
				result.Diffs = append(result.Diffs, fmt.Sprintf("env %q: declared=%q actual=%q", k, declaredVal, actualVal))
			}
		}

		if len(result.Diffs) > 0 {
			result.Drifted = true
		}

		results = append(results, result)
	}

	return results, nil
}

// Summary returns a human-readable summary of all drift results.
func Summary(results []DriftResult) string {
	var sb strings.Builder
	for _, r := range results {
		if r.Drifted {
			sb.WriteString(fmt.Sprintf("[DRIFT] %s\n", r.ServiceName))
			for _, d := range r.Diffs {
				sb.WriteString(fmt.Sprintf("  - %s\n", d))
			}
		} else {
			sb.WriteString(fmt.Sprintf("[OK]    %s\n", r.ServiceName))
		}
	}
	return sb.String()
}
