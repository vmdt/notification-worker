package cron

import (
	"fmt"
	"sync"

	"github.com/robfig/cron/v3"
)

type CronManager struct {
	c         *cron.Cron
	tasks     map[string]cron.EntryID
	taskFuncs map[string]func()
	mu        sync.Mutex
}

var (
	defaultManager *CronManager
	once           sync.Once
)

func NewCronManager() *CronManager {
	c := cron.New()
	c.Start()
	return &CronManager{
		c:         c,
		tasks:     make(map[string]cron.EntryID),
		taskFuncs: make(map[string]func()),
	}
}

func GetDefaultManager() *CronManager {
	once.Do(func() {
		defaultManager = NewCronManager()
	})
	return defaultManager
}

// AddTask adds a scheduled task with a unique name
func (m *CronManager) AddTask(name string, schedule string, fn func()) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.taskFuncs[name]; exists {
		return fmt.Errorf("task `%s` already exists", name)
	}

	entryID, err := m.c.AddFunc(schedule, fn)
	if err != nil {
		return fmt.Errorf("failed to add task `%s`: %w", name, err)
	}

	m.taskFuncs[name] = fn
	m.tasks[name] = entryID
	fmt.Println("Added task:", name)
	return nil
}

func (m *CronManager) RemoveTask(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	entryID, exists := m.tasks[name]
	if !exists {
		return fmt.Errorf("task `%s` does not exist", name)
	}

	m.c.Remove(entryID)
	delete(m.tasks, name)
	delete(m.taskFuncs, name)
	fmt.Println("Removed task:", name)
	return nil
}

func (m *CronManager) ListTasks() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	names := make([]string, 0, len(m.taskFuncs))
	for name := range m.taskFuncs {
		names = append(names, name)
	}
	return names
}
