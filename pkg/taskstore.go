// The "taskstore" package provides a simple in-memory datastore for tasks
// it uses mutex from the package "sync" to allow concurrent access

package taskstore

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	Id   int       `json:"id"`
	Text string    `json:"text"`
	Tags []string  `json:"tags"`
	Due  time.Time `json:"due"`
}

type TaskStore struct {
	mu     sync.Mutex
	tasks  map[int]Task
	nextId int
}

func New() *TaskStore {
	ts := &TaskStore{}
	ts.tasks = make(map[int]Task)
	ts.nextId = 0
	return ts
}

// CreateTask creates a new task in the store.
func (ts *TaskStore) CreateTask(text string, tags []string, due time.Time) int {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	task := Task{
		Id:   ts.nextId,
		Text: text,
		Due:  due,
	}

	task.Tags = make([]string, len(tags))

	copy(task.Tags, tags)

	ts.tasks[ts.nextId] = task
	ts.nextId++

	return task.Id
}

// GetTask retrieves a task from the store, by id. If no such id exists, an error is returned.
func (ts *TaskStore) GetTask(id int) (Task, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if task, ok := ts.tasks[id]; ok {
		return task, nil
	}

	return Task{}, fmt.Errorf("task with id=%d does not exist", id)

}

// GetAllTasks returns all the tasks in the store, in arbitrary order.
func (ts *TaskStore) GetAllTasks() []Task {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	allTasks := make([]Task, 0, len(ts.tasks))

	for _, task := range ts.tasks {
		allTasks = append(allTasks, task)
	}

	return allTasks

}

// GetTasksByTag returns all the tasks that have the given tag, in arbitrary order.
func (ts *TaskStore) GetTasksByTag(tag string) []Task {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	var tasks []Task

	for _, task := range ts.tasks {

		for _, taskTag := range task.Tags {

			if taskTag == tag {
				tasks = append(tasks, task)
				break

			}

		}
	}
	return tasks
}

// GetTasksByDueDate returns all the tasks that have the given due date, in arbitrary order.
func (ts *TaskStore) GetTasksByDueDate(year int, month time.Month, day int) []Task {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	var tasks []Task

	for _, task := range ts.tasks {
		y, m, d := task.Due.Date()
		if y == year && m == month && d == day {
			tasks = append(tasks, task)
		}
	}

	return tasks
}

// DeleteTask deletes the task with the given id. If no such id exists, an error is returned.
func (ts *TaskStore) DeleteTask(id int) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if _, ok := ts.tasks[id]; ok {
		delete(ts.tasks, id)
		return nil
	}

	return fmt.Errorf("task with id=%d not found", id)
}

// DeleteAllTasks deletes all tasks in the store.
func (ts *TaskStore) DeleteAllTasks() error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.tasks = make(map[int]Task)
	return nil

}
