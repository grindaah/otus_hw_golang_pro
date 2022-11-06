package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type ErrCounter struct {
	sync.Mutex
	i int
}

func (ec *ErrCounter) Inc() {
	ec.Lock()
	ec.i++
	ec.Unlock()
}

//func (ec *ErrCounter) Get() int {
//	return ec.i
//}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	//var wg1 sync.WaitGroup
	//var wgErrors sync.WaitGroup
	chanErrorsExceeded := make(chan struct{})
	chanErrors := make(chan error)
	chanTasks := make(chan Task)
	chanDone := make(chan struct{})
	var errCounter ErrCounter
	var wg sync.WaitGroup
	wg.Add(n)

	Producer := func(tasksChan chan<- Task) {
		for i := range tasks {
			tasksChan <- tasks[i]
		}
		close(tasksChan)
		return
	}

	ErrChecker := func() {
		for {
			select {
			case e := <-chanErrors:
				fmt.Println("error received", e.Error())
				errCounter.Inc()
				if errCounter.i > m {
					chanErrorsExceeded <- struct{}{}
					chanDone <- struct{}{}
					return
				}
			case <-chanDone:
				return
			}
		}
	}

	//Consumer
	RunWorker := func(wgMain *sync.WaitGroup) {
		for {
			select {
			case <-chanErrorsExceeded:
				fmt.Println("exit by errors")
				wgMain.Done()
				return
			case t, ok := <-chanTasks:
				if !ok {
					wgMain.Done()
					chanDone <- struct{}{}
				}
				err := t()
				if err != nil {
					chanErrors <- err
				}
			}
		}

	}

	go ErrChecker()
	go Producer(chanTasks)
	for i := 0; i < n; i++ {
		go RunWorker(&wg)
	}

	wg.Wait()
	for i := 0; i < n+1; i++ {
		chanDone <- struct{}{}
	}

	return nil
}
