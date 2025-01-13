package main

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type Task struct {
	Id     int
	Workload int
	Callback func() error
}

type TaskQueue struct {
	tasks      []Task
	mutex      sync.Mutex
	taskChan   chan Task
	quitChan   chan struct{}
	wg         sync.WaitGroup
	maxWorkers int
}

func NewTaskQueue(maxWorkers int) *TaskQueue {
	return &TaskQueue{
		taskChan:  make(chan Task),
		quitChan:  make(chan struct{}),
		maxWorkers: maxWorkers,
	}
}

func (q *TaskQueue) Start() {
	for i := 0; i < q.maxWorkers; i++ {
		go q.worker()
	}
}

func (q *TaskQueue) Stop() {
	close(q.quitChan)
	q.wg.Wait()
}

func (q *TaskQueue) Enqueue(task Task) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.tasks = append(q.tasks, task)

	// Sort the tasks based on workload priority
	sort.Slice(q.tasks, func(i, j int) bool {
		return q.tasks[i].Workload < q.tasks[j].Workload
	})

	// Send the first task in the queue to the worker
	if len(q.tasks) > 0 {
		select {
		case q.taskChan <- q.tasks[0]:
			q.tasks = q.tasks[1:]
		default:
		}
	}
}

func (q *TaskQueue) worker() {
	for {
		select {
		case task := <-q.taskChan:
			q.wg.Add(1)
			defer q.wg.Done()

			if err := task.Callback(); err != nil {
				fmt.Printf("Task %d failed: %v\n", task.Id, err)
			} else {
				fmt.Printf("Task %d completed\n", task.Id)
			}

			// Enqueue the next task in the queue if available
			q.mutex.Lock()
			if len(q.tasks) > 0 {
				select {
				case q.taskChan <- q.tasks[0]:
					q.tasks = q.tasks[1:]
				default:
				}
			}
			q.mutex.Unlock()

		case <-q.quitChan:
			return
		}
	}
}

func main() {
	// Create a new task queue with 3 workers
	queue := NewTaskQueue(3)
	queue.Start()
	defer queue.Stop()

	// Simulate enqueueing tasks with different workloads
	for i := 1; i <= 10; i++ {
		workload := rand.Intn(10) + 1
		task := Task{
			Id:       i,
			Workload: workload,
			Callback: func() error {
				time.Sleep(time.Duration(workload) * time.Second)
				// Simulate an error condition for task 3
				if i == 3 {
					return fmt.Errorf("task %d failed unexpectedly", i)
				}
				return nil
			},
		}
		fmt.Printf("Enqueuing task %d with workload %d\n", task.Id, task.Workload)
		queue.Enqueue(task)
	}

	// Wait for all tasks to complete
	fmt.Println("Waiting for all tasks to complete...")
	queue.wg.Wait()
	fmt.Println("All tasks completed.")