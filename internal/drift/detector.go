package drift

// Detail describes a single drift discrepancy for a service field.
type Detail struct {
	Service  string
	Field    string
	Declared string
	Live     string
}

// Summary holds the aggregated results of a drift detection run.
type Summary struct {
	Total   int
	Drifted int
	Missing int
	Details []Detail
}

// ServiceSpec is the declared desired state of a service.
type ServiceSpec struct {
	Name  string
	Image string
	Env   map[string]string
}

// LiveState represents the observed running state of a service.
type LiveState struct {
	Name  string
	Image string
	Env   map[string]string
}

// Detect compares declared service specs against live states and returns a Summary.
func Detect(declared []ServiceSpec, live map[string]LiveState) Summary {
	summary := Summary{Total: len(declared)}

	for _, spec := range declared {
		state, found := live[spec.Name]
		if !found {
			summary.Missing++
			summary.Drifted++
			summary.Details = append(summary.Details, Detail{
				Service:  spec.Name,
				Field:    "existence",
				Declared: "present",
				Live:     "missing",
			})
			continue
		}

		serviceDrifted := false

		if spec.Image != state.Image {
			serviceDrifted = true
			summary.Details = append(summary.Details, Detail{
				Service:  spec.Name,
				Field:    "image",
				Declared: spec.Image,
				Live:     state.Image,
			})
		}

		for k, declaredVal := range spec.Env {
			liveVal, ok := state.Env[k]
			if !ok || liveVal != declaredVal {
				serviceDrifted = true
				summary.Details = append(summary.Details, Detail{
					Service:  spec.Name,
					Field:    "env:" + k,
					Declared: declaredVal,
					Live:     liveVal,
				})
			}
		}

		if serviceDrifted {
			summary.Drifted++
		}
	}

	return summary
}
