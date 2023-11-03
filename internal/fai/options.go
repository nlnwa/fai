package fai

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/nlnwa/fai/internal/log"
)

type options struct {
	sourceDir        string
	validTargetDir   string
	invalidTargetDir string
	tmpDir           string
	concurrency      int
	sleep            time.Duration
	globPattern      string
	logger           *slog.Logger
}

func defaultOptions() *options {
	return &options{
		sourceDir:        "",
		validTargetDir:   "",
		invalidTargetDir: "",
		tmpDir:           "",
		concurrency:      runtime.NumCPU(),
		sleep:            1 * time.Second,
		globPattern:      "*",
		logger:           log.Noop(),
	}
}

func validateOptions(options ...Option) (*options, error) {
	o := defaultOptions()
	for _, opt := range options {
		opt(o)
	}

	if o.concurrency < 1 {
		return nil, fmt.Errorf("concurrency must be greater than 0")
	}

	var err error

	if (o.sourceDir == o.validTargetDir || o.sourceDir == o.invalidTargetDir) && (strings.HasSuffix(o.globPattern, "*") || strings.HasSuffix(o.globPattern, checksumFileSuffix)) {
		return nil, fmt.Errorf("source and target directories cannot be the same when glob pattern is a wildcard or glob pattern ends with %s", checksumFileSuffix)
	}

	// make sure source, temp and target directories are absolute paths
	o.sourceDir, err = filepath.Abs(o.sourceDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path of source directory: %w", err)
	}
	if info, err := os.Stat(o.sourceDir); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("source directory does not exist: %w", err)
	} else if !info.IsDir() {
		return nil, fmt.Errorf("source directory is not a directory: %w", err)
	}

	o.validTargetDir, err = filepath.Abs(o.validTargetDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path of valid target directory: %w", err)
	}
	if info, err := os.Stat(o.validTargetDir); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("valid target directory does not exist: %w", err)
	} else if !info.IsDir() {
		return nil, fmt.Errorf("valid target directory is not a directory: %w", err)
	}
	o.invalidTargetDir, err = filepath.Abs(o.invalidTargetDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path of invalid target directory: %w", err)
	}
	if info, err := os.Stat(o.invalidTargetDir); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("invalid target directory does not exist: %w", err)
	} else if !info.IsDir() {
		return nil, fmt.Errorf("invalid target directory is not a directory: %w", err)
	}
	o.tmpDir, err = filepath.Abs(o.tmpDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path of tmp directory: %w", err)
	}
	if info, err := os.Stat(o.tmpDir); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("tmp directory does not exist: %w", err)
	} else if !info.IsDir() {
		return nil, fmt.Errorf("tmp directory is not a directory: %w", err)
	}

	o.globPattern = filepath.Join(o.sourceDir, o.globPattern)

	// Test glob pattern is valid
	_, err = filepath.Glob(o.globPattern)
	if err != nil {
		return nil, fmt.Errorf("invalid glob pattern: %w", err)
	}

	if o.logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	return o, nil
}

type Option func(opts *options)

func WithSourceDir(dir string) Option {
	return func(opts *options) {
		opts.sourceDir = dir
	}
}

func WithValidTargetDir(dir string) Option {
	return func(opts *options) {
		opts.validTargetDir = dir
	}
}

func WithInvalidTargetDir(dir string) Option {
	return func(opts *options) {
		opts.invalidTargetDir = dir
	}
}

func WithTmpDir(dir string) Option {
	return func(opts *options) {
		opts.tmpDir = dir
	}
}

func WithConcurrency(concurrency int) Option {
	return func(opts *options) {
		opts.concurrency = concurrency
	}
}

func WithSleep(sleep time.Duration) Option {
	return func(opts *options) {
		opts.sleep = sleep
	}
}

func WithGlobPattern(globPattern string) Option {
	return func(opts *options) {
		opts.globPattern = globPattern
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}
