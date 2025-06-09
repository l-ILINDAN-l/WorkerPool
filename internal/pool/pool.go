package pool

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type Pool struct {
	workers      map[int]*Worker
	jobs         chan string
	addWorker    chan struct{}
	removeWorker chan struct{}
	shutdown     chan struct{}
	mu           sync.Mutex
	nextWorkerID int
}

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

func (p *Pool) AddWorker() {
	p.addWorker <- struct{}{}
}

func (p *Pool) addNewWorker() {
	p.mu.Lock()
	defer p.mu.Unlock()

	worker := NewWorker(p.nextWorkerID, p.jobs)

	p.workers[worker.id] = worker
	p.nextWorkerID++

	go worker.Start()
}

func (p *Pool) RemoveWorker() {
	p.removeWorker <- struct{}{}
}

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

func (p *Pool) Shutdown() {
	close(p.shutdown)
}

func (p *Pool) SubmitJob(job string) {
	p.jobs <- job
}
