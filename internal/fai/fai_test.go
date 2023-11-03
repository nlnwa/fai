package fai

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

// createTestFile creates a test file in the given directory
func createTestFile(t *testing.T, dir string) *os.File {
	t.Helper()
	f, _ := os.CreateTemp(dir, "testfile")
	return f
}

func TestRun(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()

	testFiles := []struct {
		file        *os.File
		expectedDir string
		isValid     bool
	}{
		{
			file:        createTestFile(t, sourceDir),
			expectedDir: targetDir,
		},
		{
			file:        createTestFile(t, sourceDir),
			expectedDir: targetDir,
		},
	}

	fai, err := New(
		WithConcurrency(len(testFiles)),
		WithSleep(0),
		WithSourceDir(sourceDir),
		WithValidTargetDir(targetDir),
		WithInvalidTargetDir(targetDir),
		WithTmpDir(t.TempDir()),
	)
	if err != nil {
		t.Fatalf("failed to create fai: %v", err)
	}

	// run fai
	fai.Run(context.Background())

	for _, testFile := range testFiles {
		// check that test file has been moved
		_, err = os.Stat(filepath.Join(testFile.expectedDir, filepath.Base(testFile.file.Name())))
		if err != nil {
			t.Errorf("failed to stat test file '%v'", err)
		}
		// check that checksum file has been created and moved
		_, err = os.Stat(filepath.Join(testFile.expectedDir, filepath.Base(testFile.file.Name())+".md5"))
		if err != nil {
			t.Errorf("failed to stat checksum file: %v", err)
		}
	}
}
