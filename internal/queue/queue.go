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

package queue

import (
	"sync"
)

// workQueue is a queue of jobs that are executed by a number of workers.
type workQueue struct {
	queue chan string
	wg    sync.WaitGroup
	m     sync.Mutex
	hm    map[string]struct{}
}

// NewWorkQueue creates a new work queue and adds the given number of workers.
func NewWorkQueue(execute func(string), concurrency int) *workQueue {
	iw := &workQueue{
		queue: make(chan string, concurrency),
		hm:    make(map[string]struct{}, concurrency),
	}

	for range concurrency {
		iw.wg.Add(1)
		go func() {
			defer iw.wg.Done()
			for job := range iw.queue {
				execute(job)
				iw.m.Lock()
				delete(iw.hm, job)
				iw.m.Unlock()
			}
		}()
	}

	return iw
}

// CloseAndWait closes the queue and waits for all workers to complete.
func (iw *workQueue) CloseAndWait() {
	// close queue
	close(iw.queue)
	// and wait for queue to be drained
	iw.wg.Wait()
}

// Add adds a job to the queue.
// If the job is already in the queue, it will be ignored.
// If the queue is full, it will block until there is room.
// If the queue is closed, it will panic.
func (iw *workQueue) Add(job string) {
	iw.m.Lock()
	// check if job is already in queue
	if _, ok := iw.hm[job]; ok {
		iw.m.Unlock()
		return
	}
	// add job to queue
	iw.hm[job] = struct{}{}
	iw.m.Unlock()

	iw.queue <- job
}
