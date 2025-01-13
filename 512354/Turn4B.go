package main

import (
	"fmt"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type Task struct {
	priority int
	callback func() error
}

type TaskQueue struct {
	tasks     []Task
	mu        sync.Mutex
	active    uint64 // Counter for active goroutines
	maxActive uint64
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{tasks: []Task{}}
}

func (q *TaskQueue) Add(task Task) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.tasks = append(q.tasks, task)
	sort.Slice(q.tasks, func(i, j int) bool {
		return q.tasks[i].priority < q.tasks[j].priority
	})
}

func (q *TaskQueue) GetNext() Task {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.tasks) == 0 {
		return Task{}
	}
	task := q.tasks[0]
	q.tasks = q.tasks[1:]
	return task
}

func (q *TaskQueue) Worker(done chan<- bool) {
	for {
		task := q.GetNext()
		if task.callback == nil {
			break
		}
		atomic.AddUint64(&q.active, 1)
		if err := task.callback(); err != nil {
			fmt.Println("Task failed:", err)
		} else {
			fmt.Println("Task completed successfully.")
		}
		atomic.AddUint64(&q.active, ^uint64(0))
	}
	done <- true
}

func (q *TaskQueue) Monitor() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			active := atomic.LoadUint64(&q.active)
			fmt.Printf("Active goroutines: %d, Max active goroutines: %d\n", active, q.maxActive)
			if active > q.maxActive {
				q.maxActive = active
			}
		}
	}
}

func main() {
	queue := NewTaskQueue()

	go queue.Monitor()

	numWorkers := runtime.NumCPU()
	fmt.Println("Starting", numWorkers, "workers...")
	done := make(chan bool, numWorkers)
	for i := 0; i < numWorkers; i++ {
		go queue.Worker(done)
	}

	// Add tasks to the queue...

	// Wait for all workers to finish
	for i := 0; i < numWorkers; i++ {
		<-done
	}
}
