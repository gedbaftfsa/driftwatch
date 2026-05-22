package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestMain_BuildAndHelp verifies the binary compiles and responds to --help.
func TestMain_BuildAndHelp(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping build test in short mode")
	}

	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "driftwatch")

	build := exec.Command("go", "build", "-o", binPath, ".")
	build.Dir = "."
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}

	cmd := exec.Command(binPath, "--help")
	out, err := cmd.CombinedOutput()
	// flag package exits with code 2 for --help
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if !ok || exitErr.ExitCode() != 2 {
			t.Fatalf("unexpected error running --help: %v\n%s", err, out)
		}
	}

	for _, flag := range []string{"-config", "-format", "-snapshot-dir", "-save-snapshot"} {
		if !containsString(string(out), flag) {
			t.Errorf("expected --help output to mention flag %q", flag)
		}
	}
}

// TestMain_MissingConfig verifies exit code 1 when config file is absent.
func TestMain_MissingConfig(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping build test in short mode")
	}

	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "driftwatch")

	build := exec.Command("go", "build", "-o", binPath, ".")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}

	missingCfg := filepath.Join(tmpDir, "nonexistent.yaml")
	cmd := exec.Command(binPath, "-config", missingCfg)
	cmd.Env = append(os.Environ())

	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit when config is missing, got nil")
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected *exec.ExitError, got %T", err)
	}
	if exitErr.ExitCode() != 1 {
		t.Errorf("expected exit code 1, got %d", exitErr.ExitCode())
	}
}

func containsString(haystack, needle string) bool {
	return len(haystack) >= len(needle) &&
		(haystack == needle ||
			len(haystack) > 0 && stringContains(haystack, needle))
}

func stringContains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
