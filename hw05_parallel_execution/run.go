package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

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
func Run(tasks []Task, n, m int) (err error) {
	var done = make(chan bool)
	var tasksChan = make(chan Task, 1)
	var errs = make(chan error, 1)
	var closeWorker = make(chan Counter)
	var wg sync.WaitGroup
	var errCounter Counter

	wg.Add(n)

	produce := func() {
		defer func() {
			//close(tasksChan)
			done <- true
			closeWorker <- Counter{i: n}
			fmt.Println("closing produce")
		}()
		for i := 0; i < len(tasks); {
			select {
			case <-errs:
				errCounter.Inc()
				fmt.Println("errCounter=", errCounter.Get(), ",m=", m)
				if errCounter.Get() >= m {
					fmt.Println("exceeded errors")
					err = ErrErrorsLimitExceeded
					return
				}
			default:
				select {
				case tasksChan <- tasks[i]:
					fmt.Println("sending task", i)
					i++
					continue
				default:
					//fmt.Println("default select")
					//time.Sleep(time.Millisecond * 100)
				}
			}
		}
		fmt.Println("exited from cycle")
	}

	consume := func(waitGroup *sync.WaitGroup, number int) {
		defer func() {
			fmt.Println("done...", number)
			waitGroup.Done()
		}()
		for {
			select {
			case t := <-tasksChan:
				err := t()
				if err != nil {
					fmt.Println("received err", err.Error())
					select {
					case errs <- err:
					default:
						fmt.Println("dropped")
					}

				} else {
					fmt.Println("task completed")
				}
			default:
				select {
				case remained := <-closeWorker:
					fmt.Println("closing worker", number, "remained", remained.i)
					//propagate to others
					remained.Dec()
					if remained.i > 0 {
						closeWorker <- remained
					}
					return
				default:
					//time.Sleep(time.Millisecond * 50)
				}
			}
			//t, haveMore := <-tasksChan
			//if haveMore {
			//	err := t()
			//	if err != nil {
			//		fmt.Println("received err", err.Error())
			//		errs <- err
			//	} else {
			//		fmt.Println("task completed")
			//	}
			//} else {
			//	return
			//}
		}
	}

	go produce()
	for i := 0; i < n; i++ {
		go consume(&wg, i)
	}
	<-done
	wg.Wait()

	return err
}
