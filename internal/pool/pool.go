package pool

import (
	"github.com/sirupsen/logrus"
	"sync"
)

// Pool manages a collection of workers, dispatches jobs,
// and handles dynamic scaling of the worker count.
type Pool struct {
	workers      map[int]*Worker
	jobs         chan string
	addWorker    chan struct{}
	removeWorker chan struct{}
	shutdown     chan struct{}
	mu           sync.Mutex
	nextWorkerID int
}

// NewPool creates and initializes a new worker pool with a specified
// number of initial workers.
func NewPool(initialWorkers int) *Pool {
	p := &Pool{
		workers:      make(map[int]*Worker),
		jobs:         make(chan string),
		addWorker:    make(chan struct{}),
		removeWorker: make(chan struct{}),
		shutdown:     make(chan struct{}),
		nextWorkerID: 1,
	}

	logrus.Infof("Creating a pool with %d initial workers", initialWorkers)
	for i := 0; i < initialWorkers; i++ {
		p.addNewWorker()
	}

	return p
}

// AddWorker sends a signal to the pool to add a new worker.
func (p *Pool) AddWorker() {
	p.addWorker <- struct{}{}
}

// addNewWorker creates a new worker, adds it to the pool, and starts it.
// This is an internal method and is not thread-safe without external locking.
func (p *Pool) addNewWorker() {
	p.mu.Lock()
	defer p.mu.Unlock()

	worker := NewWorker(p.nextWorkerID, p.jobs)

	p.workers[worker.id] = worker
	p.nextWorkerID++

	go worker.Start()
}

// RemoveWorker sends a signal to the pool to remove one worker.
func (p *Pool) RemoveWorker() {
	p.removeWorker <- struct{}{}
}

// popWorker stops and removes one worker from the pool.
// It targets the worker with the highest ID for simplicity.
func (p *Pool) popWorker() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.workers) == 0 {
		logrus.Info("No workers left")
		return
	}

	maxWorkerID := -1
	for id := range p.workers {
		if id > maxWorkerID {
			maxWorkerID = id
		}
	}

	if worker, ok := p.workers[maxWorkerID]; ok {
		worker.Stop()
		delete(p.workers, maxWorkerID)
	}
}

// Start runs the pool's main management loop in a separate goroutine.
// This loop listens for control signals to scale or shut down the pool
func (p *Pool) Start() {
	go func() {
		logrus.Info("Starting pool")
		for {
			select {
			case <-p.addWorker:
				p.addNewWorker()
			case <-p.removeWorker:
				p.popWorker()
			case <-p.shutdown:
				p.stopAllWorkers()
				logrus.Info("Pool shutting down")
				return
			}
		}
	}()
}

// stopAllWorkers stops all active workers and closes the jobs channel.
func (p *Pool) stopAllWorkers() {
	p.mu.Lock()
	defer p.mu.Unlock()

	logrus.Info("Stopping all workers")
	for id, worker := range p.workers {
		worker.Stop()
		delete(p.workers, id)
	}
	logrus.Info("All workers are stopped")

	close(p.jobs)
}

// Shutdown initiates a graceful shutdown of the entire worker pool.
func (p *Pool) Shutdown() {
	close(p.shutdown)
}

// SubmitJob sends a new job to the jobs channel for a worker to process.
func (p *Pool) SubmitJob(job string) {
	p.jobs <- job
}
