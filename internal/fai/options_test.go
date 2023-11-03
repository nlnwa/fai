package fai

import (
	"os"
	"testing"
)

func TestValidateOptions(t *testing.T) {
	testDir := t.TempDir()
	testFile, err := os.CreateTemp(testDir, "test")
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Test glob pattern is a wildcard and source and valid target directory are the same
	_, err = validateOptions(
		WithSourceDir(testDir),
		WithValidTargetDir(testDir),
		WithGlobPattern("*"),
	)
	if err == nil {
		t.Error("expected error when source and valid target directory are the same and glob pattern is '*'")
	}

	// Test glob pattern is a wildcard and source and invalid target directory are the same
	_, err = validateOptions(
		WithSourceDir(testDir),
		WithInvalidTargetDir(testDir),
		WithGlobPattern("*"),
	)
	if err == nil {
		t.Error("expected error when source and invalid target directory are the same and glob pattern is '*'")
	}

	// Test glob pattern ends with checksum file suffix and source and target directory are the same
	_, err = validateOptions(
		WithSourceDir(testDir),
		WithValidTargetDir(testDir),
		WithGlobPattern("*."+checksumFileSuffix),
	)
	if err == nil {
		t.Errorf("expected error when source and target directory are the same and glob pattern ends with '%s'", checksumFileSuffix)
	}

	// Test valid target directory is not a directory
	_, err = validateOptions(
		WithSourceDir(testDir),
		WithValidTargetDir(testFile.Name()),
	)
	if err == nil {
		t.Error("expected error when valid target directory is not a directory")
	}

	// Test concurrency cannot be zero
	_, err = validateOptions(
		WithConcurrency(0),
	)
	if err == nil {
		t.Error("expected error when concurrency is zero")
	}

	// Test glob pattern is invalid
	_, err = validateOptions(
		WithGlobPattern("["),
	)
	if err == nil {
		t.Error("expected error when glob pattern is invalid")
	}

	// Test logger cannot be nil
	_, err = validateOptions(
		WithLogger(nil),
	)
	if err == nil {
		t.Error("expected error when logger is nil")
	}
}
