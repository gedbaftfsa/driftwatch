package diff

import (
	"strings"
	"testing"
)

func TestCompareFields_NoDiff(t *testing.T) {
	declared := map[string]string{"image": "nginx:1.25", "replicas": "3"}
	live := map[string]string{"image": "nginx:1.25", "replicas": "3"}

	diffs := CompareFields(declared, live)
	if len(diffs) != 0 {
		t.Fatalf("expected no diffs, got %d", len(diffs))
	}
}

func TestCompareFields_ValueChanged(t *testing.T) {
	declared := map[string]string{"image": "nginx:1.25"}
	live := map[string]string{"image": "nginx:1.26"}

	diffs := CompareFields(declared, live)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Field != "image" || diffs[0].Declared != "nginx:1.25" || diffs[0].Live != "nginx:1.26" {
		t.Errorf("unexpected diff: %+v", diffs[0])
	}
}

func TestCompareFields_MissingInLive(t *testing.T) {
	declared := map[string]string{"image": "nginx:1.25", "port": "8080"}
	live := map[string]string{"image": "nginx:1.25"}

	diffs := CompareFields(declared, live)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Field != "port" || diffs[0].Live != "<missing>" {
		t.Errorf("unexpected diff: %+v", diffs[0])
	}
}

func TestCompareFields_UntrackedInLive(t *testing.T) {
	declared := map[string]string{"image": "nginx:1.25"}
	live := map[string]string{"image": "nginx:1.25", "extra": "value"}

	diffs := CompareFields(declared, live)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Field != "extra" || diffs[0].Declared != "<untracked>" {
		t.Errorf("unexpected diff: %+v", diffs[0])
	}
}

func TestFieldDiff_String(t *testing.T) {
	f := FieldDiff{Field: "image", Declared: "nginx:1.25", Live: "nginx:1.26"}
	s := f.String()
	if !strings.Contains(s, "image") || !strings.Contains(s, "nginx:1.25") || !strings.Contains(s, "nginx:1.26") {
		t.Errorf("unexpected string output: %s", s)
	}
}

func TestServiceDiff_HasDiff(t *testing.T) {
	empty := ServiceDiff{ServiceName: "svc", Fields: nil}
	if empty.HasDiff() {
		t.Error("expected no diff for empty fields")
	}

	withDiff := ServiceDiff{
		ServiceName: "svc",
		Fields:      []FieldDiff{{Field: "image", Declared: "a", Live: "b"}},
	}
	if !withDiff.HasDiff() {
		t.Error("expected diff to be detected")
	}
}

func TestServiceDiff_Summary_NoDiff(t *testing.T) {
	sd := ServiceDiff{ServiceName: "api", Fields: nil}
	s := sd.Summary()
	if !strings.Contains(s, "no diff") {
		t.Errorf("expected 'no diff' in summary, got: %s", s)
	}
}

func TestServiceDiff_Summary_WithDiff(t *testing.T) {
	sd := ServiceDiff{
		ServiceName: "api",
		Fields: []FieldDiff{
			{Field: "image", Declared: "nginx:1.25", Live: "nginx:1.26"},
		},
	}
	s := sd.Summary()
	if !strings.Contains(s, "api") || !strings.Contains(s, "image") {
		t.Errorf("unexpected summary output: %s", s)
	}
}
