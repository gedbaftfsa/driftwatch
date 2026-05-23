package baseline_test

import (
	"os"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/baseline"
)

func makeBaseline() *baseline.Baseline {
	return &baseline.Baseline{
		CreatedAt: time.Now().UTC(),
		Version:   "v1",
		Services: map[string]baseline.Service{
			"api": {
				Name:        "api",
				Image:       "api:1.0",
				Replicas:    2,
				Environment: map[string]string{"LOG_LEVEL": "info"},
			},
		},
	}
}

func TestStore_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	store, err := baseline.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	b := makeBaseline()
	if err := store.Save("prod", b); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := store.Load("prod")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Version != b.Version {
		t.Errorf("version: got %q want %q", loaded.Version, b.Version)
	}
	if _, ok := loaded.Services["api"]; !ok {
		t.Error("expected service 'api' in loaded baseline")
	}
}

func TestStore_LoadMissing(t *testing.T) {
	dir := t.TempDir()
	store, _ := baseline.NewStore(dir)
	_, err := store.Load("nonexistent")
	if err == nil {
		t.Fatal("expected error loading missing baseline")
	}
}

func TestStore_List(t *testing.T) {
	dir := t.TempDir()
	store, _ := baseline.NewStore(dir)
	b := makeBaseline()
	_ = store.Save("prod", b)
	_ = store.Save("staging", b)

	labels, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(labels) != 2 {
		t.Errorf("expected 2 labels, got %d", len(labels))
	}
}

func TestNewStore_CreatesDir(t *testing.T) {
	dir := t.TempDir() + "/nested/baselines"
	_, err := baseline.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}
