package baseline

import "fmt"

// Delta describes a single difference between a baseline and current state.
type Delta struct {
	Service string
	Field   string
	Old     string
	New     string
}

// String returns a human-readable description of the delta.
func (d Delta) String() string {
	return fmt.Sprintf("[%s] %s: %q -> %q", d.Service, d.Field, d.Old, d.New)
}

// Compare returns all deltas between a saved baseline and the current services.
// current is a map of service name to Service representing the live state.
func Compare(b *Baseline, current map[string]Service) []Delta {
	var deltas []Delta

	for name, base := range b.Services {
		cur, ok := current[name]
		if !ok {
			deltas = append(deltas, Delta{
				Service: name,
				Field:   "existence",
				Old:     "present",
				New:     "missing",
			})
			continue
		}

		if base.Image != cur.Image {
			deltas = append(deltas, Delta{
				Service: name, Field: "image",
				Old: base.Image, New: cur.Image,
			})
		}

		if base.Replicas != cur.Replicas {
			deltas = append(deltas, Delta{
				Service: name, Field: "replicas",
				Old: fmt.Sprintf("%d", base.Replicas),
				New: fmt.Sprintf("%d", cur.Replicas),
			})
		}

		for k, bv := range base.Environment {
			cv, exists := cur.Environment[k]
			if !exists {
				deltas = append(deltas, Delta{
					Service: name, Field: "env:" + k,
					Old: bv, New: "<unset>",
				})
			} else if bv != cv {
				deltas = append(deltas, Delta{
					Service: name, Field: "env:" + k,
					Old: bv, New: cv,
				})
			}
		}
	}

	for name := range current {
		if _, ok := b.Services[name]; !ok {
			deltas = append(deltas, Delta{
				Service: name,
				Field:   "existence",
				Old:     "missing",
				New:     "present",
			})
		}
	}

	return deltas
}

// HasDrift returns true if Compare finds any differences between the baseline
// and the current state. It is a convenience wrapper for callers that only
// need to know whether drift exists, not the full list of deltas.
func HasDrift(b *Baseline, current map[string]Service) bool {
	return len(Compare(b, current)) > 0
}
