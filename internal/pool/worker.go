package pool

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type Worker struct {
	id   int
	jobs <-chan string
	quit chan struct{}
}

func NewWorker(id int, jobs <-chan string) *Worker {
	return &Worker{
		id:   id,
		jobs: jobs,
		quit: make(chan struct{}),
	}
}

func (w *Worker) processJob(job string) {
	logrus.WithFields(logrus.Fields{
		"worker_id": w.id,
		"job":       job,
	}).Info("processing job")

	fmt.Printf("Worker %d processing job %s\n", w.id, job)

	logrus.WithFields(logrus.Fields{"worker_id": w.id, "job": job}).Info("Success processed job")
}

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

func (w *Worker) Stop() {
	close(w.quit)
}
