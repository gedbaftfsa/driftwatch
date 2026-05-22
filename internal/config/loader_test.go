package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/driftwatch/driftwatch/internal/config"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "driftwatch.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}
	return p
}

func TestLoadConfig_Valid(t *testing.T) {
	content := `
version: "1"
services:
  - name: api
    image: myrepo/api:latest
    replicas: 3
    ports: [8080]
    environment:
      LOG_LEVEL: info
    resources:
      cpu: "500m"
      memory: "256Mi"
`
	p := writeTempConfig(t, content)
	cfg, err := config.LoadConfig(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(cfg.Services))
	}
	svc := cfg.Services[0]
	if svc.Name != "api" {
		t.Errorf("expected service name 'api', got %q", svc.Name)
	}
	if svc.Replicas != 3 {
		t.Errorf("expected 3 replicas, got %d", svc.Replicas)
	}
}

func TestLoadConfig_MissingVersion(t *testing.T) {
	content := `services:
  - name: api
    image: myrepo/api:latest
`
	p := writeTempConfig(t, content)
	_, err := config.LoadConfig(p)
	if err == nil {
		t.Fatal("expected error for missing version, got nil")
	}
}

func TestLoadConfig_DuplicateServiceName(t *testing.T) {
	content := `
version: "1"
services:
  - name: api
    image: myrepo/api:v1
  - name: api
    image: myrepo/api:v2
`
	p := writeTempConfig(t, content)
	_, err := config.LoadConfig(p)
	if err == nil {
		t.Fatal("expected error for duplicate service name, got nil")
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := config.LoadConfig("/nonexistent/path/driftwatch.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
