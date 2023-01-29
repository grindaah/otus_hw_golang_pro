package hw06pipelineexecution

import (
	"fmt"
	"sync"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var wg sync.WaitGroup
	out := make(In)
	wg.Add(len(stages))

	RunAllStagesSync := func(done In) {
		for i, s := range stages {
			if i == 0 {
				out = s(in)
			} else {
				out = s(out)
			}

			wg.Done()
			fmt.Println("Done! (from pipeline)")
		}
	}
	go RunAllStagesSync(done)
	wg.Wait()

	return out //merge(done, []In{out}...)
}

func merge(done In, cs ...In) Out {
	fmt.Println("mergezzx")
	var wg sync.WaitGroup
	out := make(chan interface{})
	waitGr := make(chan bool)

	output := func(c In) {
		for n := range c {
			out <- n
		}
		fmt.Println("Done!")
		wg.Done()
	}

	fmt.Println("len(cs)=", len(cs))
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		fmt.Println("waitGr")
		waitGr <- true
	}()

	for {
		select {
		case <-waitGr:
			close(out)
			return out
		case <-done:
			return out
		default:
			continue
		}
	}
}
