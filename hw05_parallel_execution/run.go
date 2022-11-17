package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrWorkersNotPassed    = errors.New("workers n value must be passed more than 0")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
// m <= 0 treated as ignoring all errors.
func Run(tasks []Task, n, m int) (err error) {
	if n <= 0 {
		return ErrWorkersNotPassed
	}

	tasksChan := make(chan Task, 1)
	errs := make(chan error, 10)
	closeWorker := make(chan struct{})
	var wg sync.WaitGroup
	var errCounter int32
	ignoreErrors := m <= 0

	wg.Add(n)

	produce := func() {
		defer func() {
			close(closeWorker)
			fmt.Println("closing produce")
		}()
		for i := 0; i < len(tasks); {
			select {
			case <-errs:
				if !ignoreErrors {
					atomic.AddInt32(&errCounter, 1)
					if errCounter >= int32(m) {
						fmt.Println("exceeded errors")
						err = ErrErrorsLimitExceeded
						return
					}
				}
			case tasksChan <- tasks[i]:
				i++
				continue
			default:
			}
		}
		fmt.Println("exited from cycle")
	}

	go produce()
	for i := 0; i < n; i++ {
		go consume(&wg, i, tasksChan, errs, closeWorker)
	}
	wg.Wait()

	return err
}

func consume(waitGroup *sync.WaitGroup, number int, tasks <-chan Task, errs chan<- error, closeWorker <-chan struct{}) {
	defer func() {
		fmt.Println("done...", number)
		waitGroup.Done()
	}()
	for {
		select {
		case t := <-tasks:
			err := t()
			if err != nil {
				errs <- err
			}
		default:
			select {
			case <-closeWorker:
				fmt.Println("closing worker", number)
				return
			default:
			}
		}
	}
}
