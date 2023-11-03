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

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var filesize = promauto.NewHistogram(prometheus.HistogramOpts{
	Name: "file_size_bytes",
	Help: "Size of files in bytes.",
	// 1MB, 100MB, 500MB, 1GB
	Buckets: []float64{1000000, 100000000, 500000000, 1000000000},
})

var validationError = promauto.NewCounter(prometheus.CounterOpts{
	Name: "validation_error",
	Help: "Number of files with validation errors.",
})

func ValidationError() {
	validationError.Inc()
}

func Size(size int64) {
	filesize.Observe(float64(size))
}
