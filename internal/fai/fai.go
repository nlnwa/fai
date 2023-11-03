/*
 * Copyright 2023 National Library of Norway.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package fai

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/nlnwa/fai/internal/checksum"
	"github.com/nlnwa/fai/internal/metrics"
	"github.com/nlnwa/fai/internal/queue"
	"github.com/nlnwa/fai/internal/warc"
)

type fAI struct {
	sourceDir        string
	validTargetDir   string
	invalidTargetDir string
	tmpDir           string
	concurrency      int
	sleep            time.Duration
	globPattern      string
	logger           *slog.Logger
}

func New(options ...Option) (*fAI, error) {
	opts, err := validateOptions(options...)
	if err != nil {
		return nil, err
	}

	return &fAI{
		sourceDir:        opts.sourceDir,
		validTargetDir:   opts.validTargetDir,
		invalidTargetDir: opts.invalidTargetDir,
		tmpDir:           opts.tmpDir,
		concurrency:      opts.concurrency,
		sleep:            opts.sleep,
		globPattern:      opts.globPattern,
		logger:           opts.logger,
	}, nil
}

const checksumFileSuffix = ".md5"

func (f *fAI) processFile(file string) {
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		// file does not exist so skip it
		return
	}

	// calculate checksum
	md5sum, err := checksum.MD5Sum(file)
	if err != nil {
		f.logger.Error("Failed to calculate checksum", "file", file, "error", err)
		return
	}
	// create checksum file
	checksumFile, err := checksum.CreateChecksumFile(file, md5sum, checksumFileSuffix)
	if err != nil {
		f.logger.Error("Failed to create checksum file", "file", file, "error", err)
		return
	}

	// validate file
	isValid, err := warc.IsValid(file, filepath.Join(f.tmpDir, "buffer"))
	if err != nil {
		f.logger.Error("Failed to validate file", "file", file, "error", err)
		return
	}

	targetDir := f.validTargetDir
	if !isValid {
		targetDir = f.invalidTargetDir
		metrics.ValidationError()
	}

	newChecksumFile := filepath.Join(targetDir, filepath.Base(checksumFile))
	newFile := filepath.Join(targetDir, filepath.Base(file))

	// Move checksum file and file to target directory.
	//
	// The order is important because a failed move of the checksum file
	// will result in the file being checksummed again (ok). If the file
	// is moved first and the checksum file fails to move then the
	// checksum file will never be created (not ok).

	// move checksum file to new location
	err = os.Rename(checksumFile, newChecksumFile)
	if err != nil {
		f.logger.Error("Failed to move checksum file", "source", checksumFile, "target", newChecksumFile, "error", err)
		return
	}

	// move file to new location
	err = os.Rename(file, newFile)
	if err != nil {
		f.logger.Error("Failed to move file", "source", file, "target", newFile, "error", err)
		return
	}

	// get file size
	fileInfo, err := os.Stat(newFile)
	if err != nil {
		f.logger.Error("Failed to get file size", "file", newFile, "error", err)
		return
	}
	fileSizeBytes := fileInfo.Size()

	metrics.Size(fileSizeBytes)

	f.logger.Info("Processed file", "file", newFile, "size", fileSizeBytes, "md5", md5sum, "valid", isValid)
}

// Run starts the FAI.
// It will run until the context is cancelled or stop after one pass if the sleep duration is zero.
func (f *fAI) Run(ctx context.Context) {
	f.logger.Info("Starting FAI", "sourceDir", f.sourceDir, "validTargetDir", f.validTargetDir, "invalidTargetDir", f.invalidTargetDir, "concurrency", f.concurrency, "sleep", f.sleep)

	queue := queue.NewWorkQueue(f.processFile, f.concurrency)
	defer queue.CloseAndWait()

	for {
		files, _ := filepath.Glob(f.globPattern)
		for _, file := range files {
			select {
			case <-ctx.Done():
				return
			default:
				queue.Add(file)
			}
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(f.sleep):
			// do a single pass if sleep duration is zero
			if f.sleep == 0 {
				return
			}
		}
	}
}
