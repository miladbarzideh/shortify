package worker

import (
	"errors"
	"sync"

	"github.com/sirupsen/logrus"
)

type Pool interface {
	Submit(task func()) error
	StopAndWait()
	GetWorkerCount() int
	GetQueueSize() int
}

type pool struct {
	logger      *logrus.Logger
	workerCount int
	queueSize   int
	taskQueue   chan func()
	shutdown    chan bool
	wg          sync.WaitGroup
}

func NewWorkerPool(logger *logrus.Logger, workerCount int, queueSize int) Pool {
	if workerCount < 1 {
		workerCount = 1
	}

	if queueSize < 0 {
		queueSize = 0
	}

	wp := &pool{
		logger:      logger,
		workerCount: workerCount,
		queueSize:   queueSize,
		taskQueue:   make(chan func(), queueSize),
		shutdown:    make(chan bool),
	}
	go wp.dispatch()

	return wp
}

func (wp *pool) Submit(task func()) error {
	if task == nil {
		return nil
	}

	select {
	case wp.taskQueue <- task:
		wp.logger.Debug("adding task to task queue")
		return nil
	case <-wp.shutdown:
		wp.logger.Debug("pool shutting down")
		return errors.New("pool is shutting down")
	}
}

func (wp *pool) StopAndWait() {
	close(wp.shutdown)
	wp.logger.Debug("waiting for workers to shutdown")
	wp.wg.Wait()
	close(wp.taskQueue)
}

func (wp *pool) GetWorkerCount() int {
	return wp.workerCount
}

func (wp *pool) GetQueueSize() int {
	return wp.queueSize
}

func (wp *pool) dispatch() {
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

func (wp *pool) worker(id int) {
	defer wp.wg.Done()
	for {
		select {
		case task, ok := <-wp.taskQueue:
			if !ok {
				wp.logger.Debug("worker", id, " shutting down")
				return
			}
			wp.logger.Debug("worker", id, " processing task")
			task()
			wp.logger.Debug("worker", id, " finished task")
		case <-wp.shutdown:
			wp.logger.Debug("worker", id, " shutting down")
			return
		}
	}
}
