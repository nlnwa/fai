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

package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/nlnwa/fai/internal/fai"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	sourceDir := ""
	validTargetDir := ""
	invalidTargetDir := ""
	tmpDir := ""

	concurrency := runtime.NumCPU()
	sleep := 5 * time.Second
	pattern := "*"
	metricsPort := 8081

	flag.StringVar(&sourceDir, "source-dir", sourceDir, "path to source directory")
	flag.StringVar(&validTargetDir, "valid-target-dir", validTargetDir, "path to target directory where valid files and their corresponding checksum files will be moved to")
	flag.StringVar(&invalidTargetDir, "invalid-target-dir", invalidTargetDir, "path to target directory where invalid files and their corresponding checksum files will be moved to")
	flag.StringVar(&tmpDir, "tmp-dir", tmpDir, "path to directory where temporary buffer files will be stored")
	flag.IntVar(&concurrency, "concurrency", concurrency, "number of concurrent files processed")
	flag.DurationVar(&sleep, "sleep", sleep, "sleep duration between directory listings, set to 0 to only do a single run")
	flag.StringVar(&pattern, "pattern", pattern, "glob pattern used to match filenames in source directory")
	flag.IntVar(&metricsPort, "metrics-port", metricsPort, "port to expose metrics on")
	flag.Parse()

	logger := slog.Default()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		defer cancel()
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(fmt.Sprintf(":%d", metricsPort), nil)
		if err != nil {
			logger.Error("Failed to start metrics server", "error", err)
		}
	}()

	f, err := fai.New(
		fai.WithSourceDir(sourceDir),
		fai.WithValidTargetDir(validTargetDir),
		fai.WithInvalidTargetDir(invalidTargetDir),
		fai.WithTmpDir(tmpDir),
		fai.WithConcurrency(concurrency),
		fai.WithSleep(sleep),
		fai.WithGlobPattern(pattern),
		fai.WithLogger(logger),
	)
	if err != nil {
		logger.Error("", "error", err)
		os.Exit(1)
	}

	f.Run(ctx)
}
