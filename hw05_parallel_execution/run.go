package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var (
		counterError int32
		wg           sync.WaitGroup
		maxErrors    int32
	)

	maxErrors = int32(m)

	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	taskChannel := make(chan Task, len(tasks))

	for range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range taskChannel {
				if atomic.LoadInt32(&counterError) >= maxErrors {
					continue
				}

				err := t()
				if err != nil {
					atomic.AddInt32(&counterError, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		if m > 0 && atomic.LoadInt32(&counterError) >= maxErrors {
			break
		}

		taskChannel <- task
	}

	close(taskChannel)

	wg.Wait()

	if atomic.LoadInt32(&counterError) >= maxErrors {
		return ErrErrorsLimitExceeded
	}

	return nil
}
