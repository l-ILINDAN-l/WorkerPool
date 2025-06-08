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
		// TODO Create method addNewWorker()
		//p.addNewWorker()
	}

	return p
}
