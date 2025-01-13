package main

import (
	"fmt"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

// Task represents a task to be executed, with a priority and a callback function
type Task struct {
	priority int
	callback func() error
}

// TaskQueue manages a list of tasks sorted by priority
type TaskQueue struct {
	tasks         []Task
	mu            *sync.Mutex
	activeWorkers int32
}

// NewTaskQueue creates a new task queue
func NewTaskQueue() *TaskQueue {
	return &TaskQueue{tasks: []Task{}, mu: &sync.Mutex{}}
}

// Add adds a new task to the queue
func (q *TaskQueue) Add(task Task) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.tasks = append(q.tasks, task)
	sort.Slice(q.tasks, func(i, j int) bool {
		return q.tasks[i].priority < q.tasks[j].priority
	})
}

// GetNext retrieves and removes the next task with the highest priority
func (q *TaskQueue) GetNext() Task {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.tasks) == 0 {
		return Task{} // Return an empty task if the queue is empty
	}
	// Remove the task at index 0 and return it
	task := q.tasks[0]
	q.tasks = q.tasks[1:]
	return task
}

// Worker is a goroutine that processes tasks from the queue
func Worker(queue *TaskQueue) {
	for {
		task := queue.GetNext()
		if task.callback == nil {
			break // Exit if no tasks are left and no callback
		}

		// Increment the number of active workers
		atomic.AddInt32(&queue.activeWorkers, 1)
		defer atomic.AddInt32(&queue.activeWorkers, -1)

		if err := task.callback(); err != nil {
			fmt.Println("Task failed:", err)
		} else {
			fmt.Println("Task completed successfully.")
		}
	}
}

// MonitorConcurrency periodically logs the number of active goroutines and workers
func MonitorConcurrency(queue *TaskQueue, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		numGoroutines := runtime.NumGoroutine()
		numActiveWorkers := atomic.LoadInt32(&queue.activeWorkers)
		fmt.Printf("Active Goroutines: %d, Active Workers: %d\n", numGoroutines, numActiveWorkers)
	}
}

func main() {
	queue := NewTaskQueue()

	// Add tasks to the queue with different priorities
	queue.Add(Task{priority: 1, callback: func() error {
		fmt.Println("Task 1 (low priority) starting...")
		time.Sleep(time.Second)
		return nil
	}})

	queue.Add(Task{priority: 3, callback: func() error {
		fmt.Println("Task 3 (high priority) starting...")
		time.Sleep(time.Second)
		return nil
	}})

	queue.Add(Task{priority: 2, callback: func() error {
		fmt.Println("Task 2 (medium priority) starting...")
		time.Sleep(time.Second)
		return nil
	}})

	// Start worker goroutines
	const numWorkers = 5
	for i := 0; i < numWorkers; i++ {
		go Worker(queue)
	}

	// Start monitoring concurrency
	go MonitorConcurrency(queue, 2*time.Second)

	// Simulate other work or wait for the queue to be processed
	time.Sleep(5 * time.Second)
}
