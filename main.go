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
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/nlnwa/fai/internal/fai"
	"github.com/nlnwa/fai/internal/metrics"
	"github.com/nlnwa/fai/internal/queue"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	pflag.String("dir", "", "path to source directory")
	pflag.String("pattern", "*.warc.gz", "glob pattern used to match filenames in source directory")
	pflag.Int("concurrency", runtime.NumCPU(), "number of files processed concurrently")
	pflag.Duration("sleep", 5*time.Second, "sleep duration between directory listings, set to 0 to only do a single pass")
	pflag.String("s3-address", "localhost:9000", "s3 endpoint (address:port)")
	pflag.String("s3-bucket-name", "", "name of bucket to upload files to")
	pflag.String("s3-access-key-id", "", "access key ID")
	pflag.String("s3-secret-access-key", "", "secret access key")
	pflag.String("s3-token", "", "token to use for s3 authentication (optional)")
	pflag.Int("metrics-port", 8081, "port to expose metrics on")
	pflag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		logger.Error("Failed to bind flags", "error", err)
		os.Exit(1)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	sourceDir := viper.GetString("dir")
	concurrency := viper.GetInt("concurrency")
	sleep := viper.GetDuration("sleep")
	globPattern := viper.GetString("pattern")
	metricsAddr := fmt.Sprintf(":%d", viper.GetInt("metrics-port"))
	s3bucketName := viper.GetString("s3-bucket-name")
	s3address := viper.GetString("s3-address")
	s3accessKeyID := viper.GetString("s3-access-key-id")
	s3secretAccessKey := viper.GetString("s3-secret-access-key")
	s3token := viper.GetString("s3-token")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		defer cancel()
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(metricsAddr, nil)
		if errors.Is(err, http.ErrServerClosed) {
			return
		}
		if err != nil {
			logger.Error("Metrics server failed", "error", err)
		}
	}()

	s3uploader, err := fai.NewS3Uploader(
		fai.WithS3Address(s3address),
		fai.WithS3AccessKeyID(s3accessKeyID),
		fai.WithS3SecretAccessKey(s3secretAccessKey),
		fai.WithS3Token(s3token),
		fai.WithS3BucketName(s3bucketName),
	)
	if err != nil {
		logger.Error("Failed to create S3 uploader", "error", err)
		os.Exit(1)
	}
	logger.Info("S3 uploader", "bucket", s3bucketName, "address", s3address)

	worker := func(filePath string) {
		start := time.Now()
		info, err := s3uploader.Upload(ctx, filePath)
		if err != nil {
			logger.Error("Failed to upload file", "file", filePath, "error", err)
			return
		}
		metrics.Duration(time.Since(start))
		metrics.Size(info.Size)

		logger.Info("Uploaded file", "key", info.Key, "size", info.Size, "etag", info.ETag)

		err = os.Remove(filePath)
		if err != nil {
			logger.Error("Failed to remove file", "file", filePath, "error", err)
		}
	}

	queue := queue.NewWorkQueue(worker, concurrency)
	defer queue.CloseAndWait()
	logger.Info("Work queue", "concurrency", concurrency)

	f, err := fai.New(
		fai.WithSourceDir(sourceDir),
		fai.WithSleep(sleep),
		fai.WithGlobPattern(globPattern),
		fai.WithInspector(queue.Add),
	)
	if err != nil {
		logger.Error("", "error", err)
		os.Exit(1)
	}

	logger.Info("Starting FAI", "sourceDir", sourceDir, "globPattern", globPattern, "sleep", sleep.String())

	f.Run(ctx)
}
