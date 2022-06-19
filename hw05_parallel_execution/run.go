package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	taskChan, errChan := make(chan Task), make(chan error, n)
	defer close(errChan)

	exitChan := errorWatcher(errChan, m)

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() { worker(taskChan, errChan); wg.Done() }()
	}

	defer func() { close(taskChan); wg.Wait() }()

	for _, task := range tasks {
		select {
		case <-exitChan:
			return ErrErrorsLimitExceeded
		case taskChan <- task:
		}
	}

	return nil
}

func errorWatcher(errChan <-chan error, maxErr int) <-chan struct{} {
	exitChan := make(chan struct{})

	go func() {
		defer close(exitChan)

		var errCount int

		for range errChan {
			errCount++

			if maxErr > 0 && errCount == maxErr {
				return
			}
		}
	}()

	return exitChan
}

func worker(taskChan <-chan Task, errChan chan<- error) {
	for task := range taskChan {
		if err := task(); err != nil {
			errChan <- err
		}
	}
}
