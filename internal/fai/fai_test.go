package fai

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/nlnwa/fai/internal/queue"
)

func TestNew(t *testing.T) {
	// Test glob pattern is invalid
	_, err := New(
		WithGlobPattern("["),
	)
	if err == nil {
		t.Error("expected error when glob pattern is invalid")
	}
}

func TestRun(t *testing.T) {
	sourceDir := t.TempDir()
	var testFiles []*os.File

	for range 10 {
		f, _ := os.CreateTemp(sourceDir, "testfile")
		testFiles = append(testFiles, f)
	}

	worker := func(path string) {
		err := os.Remove(path)
		if err != nil {
			t.Errorf("failed to remove file: %v", err)
		}
	}

	q := queue.NewWorkQueue(worker, len(testFiles))

	fai, err := New(
		WithSleep(0), // single pass
		WithSourceDir(sourceDir),
		WithInspector(q.Add),
		WithGlobPattern("testfile*"),
	)
	if err != nil {
		t.Fatalf("failed to create fai: %v", err)
	}

	// run fai (add files to queue)
	fai.Run(context.Background())
	// close queue and wait for all workers to finish
	q.CloseAndWait()

	for _, testFile := range testFiles {
		// assert that test file is removed (as per worker function)
		_, err = os.Stat(testFile.Name())
		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("test file still exists: %s", testFile.Name())
		}
	}
}
