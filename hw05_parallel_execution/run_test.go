package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks with error less than m", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration
		var m sync.Mutex

		for i := 0; i < tasksCount; i++ {
			var err error
			if i%2 == 0 {
				err = fmt.Errorf("error from task %d", i)
			}
			tasks = append(tasks, func() error {
				taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
				time.Sleep(taskSleep)
				m.Lock()
				sumTime += taskSleep
				m.Unlock()
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 26

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)

		require.NoError(t, err)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("tasks less than workers", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration
		var m sync.Mutex

		for i := 0; i < tasksCount; i++ {
			var err error
			if i%2 == 0 {
				err = fmt.Errorf("error from task %d", i)
			}
			tasks = append(tasks, func() error {
				taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
				time.Sleep(taskSleep)
				m.Lock()
				sumTime += taskSleep
				m.Unlock()
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 14
		maxErrorsCount := 8

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)

		require.NoError(t, err)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})
}

// additional test using require.Eventually, so can add stress tests.
func TestAdditional(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("tasks less than workers (eventually)", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration
		var m sync.Mutex

		for i := 0; i < tasksCount; i++ {
			var err error
			if i%2 == 0 {
				err = fmt.Errorf("error from task %d", i)
			}
			tasks = append(tasks, func() error {
				taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
				time.Sleep(taskSleep)
				m.Lock()
				sumTime += taskSleep
				m.Unlock()
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 14
		maxErrorsCount := 8

		start := time.Now()
		_ = Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)

		require.Eventually(t, func() bool {
			return runTasksCount == int32(tasksCount)
		}, sumTime/2, time.Millisecond*20)
		require.Eventually(t, func() bool {
			return elapsedTime <= sumTime/2
		}, sumTime/2, time.Millisecond*20)
	})

	// passed with 100000 reduced for quick pass.
	t.Run("stress (eventually)", func(t *testing.T) {
		tasksCount := 1000
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration
		var m sync.Mutex

		for i := 0; i < tasksCount; i++ {
			var err error
			if i%200 == 0 {
				err = fmt.Errorf("error from task %d", i)
			}
			tasks = append(tasks, func() error {
				taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
				time.Sleep(taskSleep)
				m.Lock()
				sumTime += taskSleep
				m.Unlock()
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 25
		maxErrorsCount := 6

		start := time.Now()
		_ = Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)

		require.Eventually(t, func() bool {
			return runTasksCount == int32(tasksCount)
		}, sumTime/2, time.Millisecond*20)
		require.Eventually(t, func() bool {
			return elapsedTime <= sumTime/2
		}, sumTime/2, time.Millisecond*20)
	})

	t.Run("ignore errors", func(t *testing.T) {
		tasksCount := 1000
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration
		var m sync.Mutex

		for i := 0; i < tasksCount; i++ {
			var err error
			if i%50 == 0 {
				err = fmt.Errorf("error from task %d", i)
			}
			tasks = append(tasks, func() error {
				taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
				time.Sleep(taskSleep)
				m.Lock()
				sumTime += taskSleep
				m.Unlock()
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 25
		maxErrorsCount := 0

		start := time.Now()
		_ = Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)

		require.Eventually(t, func() bool {
			return runTasksCount == int32(tasksCount)
		}, sumTime/2, time.Millisecond*20)
		require.Eventually(t, func() bool {
			return elapsedTime <= sumTime/2
		}, sumTime/2, time.Millisecond*20)
	})
}
