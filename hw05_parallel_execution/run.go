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
	var errCounter Counter
	var wg sync.WaitGroup
	wg.Add(n)

	Producer := func(tasks []Task) <-chan Task {
		//defer wg.Done()
		chanTasks := make(chan Task, n)
		go func() {
			defer func() {
				fmt.Println("closing channel... chanTasks")
				close(chanTasks)
			}()
			for i := 0; i < len(tasks); i++ {
				select {
				case chanTasks <- tasks[i]:
					continue
				case <-chanErrorsExceeded:
					fmt.Println("receive errExceeded... closing channel")
					//for j := 0; j < n; j++ {
					chanErrorsExceeded <- "error channel"
					//}
					return
					/*default:
					time.Sleep(time.Millisecond * 200)
					fmt.Println("skipping...")*/
				}
			}
		}()
		return chanTasks
	}

	/*Dispatcher := func() {
		for {
			select {
			case <-chanErrorsExceeded:
				// errchecker will exit itself

			}
		}
	}*/

	ErrChecker := func() {
		//defer close(chanErrors)
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
			case <-chanDone:
				return
			default:
				continue
			}
		}
		return
	}

	chanTasksReady := Producer(tasks)
	//Consumer
	RunWorker := func(wgMain *sync.WaitGroup, number int) {
		defer func() {
			fmt.Println("exiting worker... ", number)
			wg.Done()
		}()
		for {
			select {
			case t, ok := <-chanTasksReady:
				fmt.Println("worker...", ok)
				if !ok {
					fmt.Println("exit by close channel", number)
					//chanDone <- struct{}{}
					return
				} else {
					err := t()
					if err != nil {
						err = fmt.Errorf("worker N %d error:%w", number, err)
						chanErrors <- err
						continue
					}
				}
			case <-chanErrorsExceeded:
				fmt.Println("worker exiting", number)
				chanErrorsExceeded <- strconv.Itoa(number)
				return
			default:
				time.Sleep(200 * time.Millisecond)
				fmt.Println("default...")
				continue
			}
		}
		return
	}

	go ErrChecker()
	//go Producer(cha)
	for i := 0; i < n; i++ {
		go RunWorker(&wg, i)
	}

	wg.Wait()
	chanDone <- "aaa"

	return err
}
