package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrWorkersNotPassed    = errors.New("workers n value must be passed more than 0")
)

type Task func() error

type Counter struct {
	sync.Mutex
	i int
}

func (c *Counter) Inc() {
	c.Lock()
	c.i++
	c.Unlock()
}

func (c *Counter) Dec() {
	c.Lock()
	c.i--
	c.Unlock()
}

func (c *Counter) Get() int {
	c.Lock()
	defer c.Unlock()
	return c.i
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
// m <= 0 treated as ignoring all errors.
func Run(tasks []Task, n, m int) (err error) {
	if n <= 0 {
		return ErrWorkersNotPassed
	}

	done := make(chan bool)
	tasksChan := make(chan Task, 1)
	errs := make(chan error, 1)
	closeWorker := make(chan *Counter)
	var wg sync.WaitGroup
	var errCounter Counter
	ignoreErrors := m <= 0

	wg.Add(n)

	produce := func() {
		defer func() {
			done <- true
			closeWorker <- &Counter{i: n}
			fmt.Println("closing produce")
		}()
		for i := 0; i < len(tasks); {
			select {
			case <-errs:
				if !ignoreErrors {
					errCounter.Inc()
					if errCounter.Get() >= m {
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
	<-done
	wg.Wait()

	return err
}

func consume(waitGroup *sync.WaitGroup, number int, tasks <-chan Task, errs chan<- error, closeWorker chan *Counter) {
	defer func() {
		fmt.Println("done...", number)
		waitGroup.Done()
	}()
	for {
		select {
		case t := <-tasks:
			err := t()
			if err != nil {
				select {
				case errs <- err:
				default:
				}
			}
		default:
			select {
			case remained := <-closeWorker:
				fmt.Println("closing worker", number, "remained", remained.i)
				// propagate to others
				remained.Dec()
				if remained.i > 0 {
					closeWorker <- remained
				}
				return
			default:
			}
		}
	}
}
