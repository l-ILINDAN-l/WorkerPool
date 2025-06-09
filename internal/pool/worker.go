package pool

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

// Worker represents a single worker that executes jobs.
// It is designed to run in its own goroutine and process jobs
// received from a shared channel.
type Worker struct {
	// id is a unique identifier for the worker.
	id int
	// jobs is a read-only channel from which the worker receives jobs.
	jobs <-chan string
	// quit is a channel to signal the worker to stop its execution loop.
	quit chan struct{}
}

// NewWorker creates and returns a new Worker instance.
// It requires a unique id for identification and the shared jobs channel.
func NewWorker(id int, jobs <-chan string) *Worker {
	return &Worker{
		id:   id,
		jobs: jobs,
		quit: make(chan struct{}),
	}
}

// processJob handles the execution of a single job.
// In this case, it logs the start and end of processing.
func (w *Worker) processJob(job string) {
	logrus.WithFields(logrus.Fields{
		"worker_id": w.id,
		"job":       job,
	}).Info("processing job")

	fmt.Printf("Worker %d processing job %s\n", w.id, job)

	logrus.WithFields(logrus.Fields{"worker_id": w.id, "job": job}).Info("Success processed job")
}

// Start runs the worker's main loop.
// The worker listens for incoming jobs or a signal to quit.
// This method should be called as a goroutine.
func (w *Worker) Start() {
	logrus.Infof("Starting worker %d\n", w.id)

	for {
		select {
		case job, ok := <-w.jobs:
			if !ok {
				logrus.Warn("Worker %d stopped", w.id)
			}

			w.processJob(job)
		case <-w.quit:
			logrus.Info("Worker %d stopped with signal quit", w.id)
			return
		}
	}

}

// Stop sends a signal to the worker's quit channel to gracefully shut it down.
func (w *Worker) Stop() {
	close(w.quit)
}
