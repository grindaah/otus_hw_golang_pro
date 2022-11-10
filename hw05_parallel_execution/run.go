package hw05parallelexecution

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
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

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) (err error) {

	chanErrorsExceeded := make(chan string, 1)
	chanErrors := make(chan error)

	chanDone := make(chan string)

	chanWorkerFinish := make(chan string, n)
	chanWorkersFinished := make(chan string)
	chanFinishProducer := make(chan string)
	chanFinishErrorChecker := make(chan string)

	var errCounter Counter
	var wg sync.WaitGroup
	wg.Add(n + 2)

	Producer := func(tasks []Task, wg *sync.WaitGroup) <-chan Task {
		//defer wg.Done()
		defer func() {
			fmt.Println("closing producer")
			wg.Done()
		}()
		chanTasks := make(chan Task, n)
		go func() {
			defer func() {
				fmt.Println("closing channel... chanTasks")
				close(chanTasks)
			}()
			for i := 0; i < len(tasks); {
				select {
				case chanTasks <- tasks[i]:
					i++
					continue
				case f := <-chanFinishProducer:
					fmt.Println("received", f)
					return
				default:
					continue
				}
			}
			return
		}()
		return chanTasks
	}

	Dispatcher := func() {
		defer func() {
			fmt.Println("closing dispatcher")
		}
		for {
			select {
			case <-chanErrorsExceeded:
				// errchecker will exit itself
				for i := 0; i < n; i++ {
					fmt.Println("sending finish signal", i)
					chanWorkerFinish <- "finish by error exceeded " + strconv.Itoa(i)
				}
				chanFinishProducer <- "from dispatcher (1)"
			case <-chanWorkersFinished:
				chanFinishProducer <- "from dispatcher (2)"
				chanFinishErrorChecker <- "from dispatcher (3)"
			case <-chanDone:
				return
			default:
				continue
			}
		}
		return
	}

	go Dispatcher()

	ErrChecker := func(wg *sync.WaitGroup) {
		defer func() {
			fmt.Println("closing errchecker")
			wg.Done()
		}()
		for {
			select {
			case e := <-chanErrors:
				fmt.Println("error received", e.Error())
				errCounter.Inc()
				if errCounter.i > m {
					fmt.Println("send close by channel condition")
					chanErrorsExceeded <- "5"
					return
				}
			case <-chanFinishErrorChecker:
				return
			default:
				continue
			}
		}
		return
	}

	chanTasksReady := Producer(tasks, &wg)
	//Consumer
	RunWorker := func(wgMain *sync.WaitGroup, number int) {
		defer func() {
			fmt.Println("exiting worker... ", number)
			wg.Done()
		}()
		for {
			select {
			case <-chanWorkerFinish:
				fmt.Println("worker exiting", number)
				//chanErrorsExceeded <- strconv.Itoa(number)
				return
			case t := <-chanTasksReady:
				fmt.Println("worker... ", number) //, haveTask)
				//if haveTask {
				err := t()
				if err != nil {
					err = fmt.Errorf("worker N %d error:%w", number, err)
					chanErrors <- err
					continue
				}

				//} else {
				//	fmt.Println("exit by close channel", number)
				//	//chanDone <- struct{}{}
				//	return
				//}
			default:
				time.Sleep(200 * time.Millisecond)
				fmt.Println("default...")
				continue
			}
		}
		return
	}

	go ErrChecker(&wg)
	//go Producer(cha)
	for i := 0; i < n; i++ {
		go RunWorker(&wg, i)
	}

	wg.Wait()
	chanDone <- "aaa"

	return err
}
