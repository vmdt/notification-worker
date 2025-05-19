package cron

import (
	"fmt"
	"sync"

	"github.com/robfig/cron/v3"
)

var (
	cronInstance *cron.Cron
	once         sync.Once
)

var taskFunc = make(map[string]func())

func NewCronInstance() *cron.Cron {
	newCron := cron.New()
	newCron.Start()
	return newCron
}

func AddTaskFunc(name string, schedule string, f func()) {
	if _, ok := taskFunc[name]; !ok {
		fmt.Println("Add a new task:", name)
		cInstance := GetInstance()
		cInstance.AddFunc(schedule, f)
		taskFunc[name] = f
	} else {
		fmt.Println("Don't add same task `" + name + "` repeatedly!")
	}
}

func GetInstance() *cron.Cron {
	once.Do(func() {
		cronInstance = NewCronInstance()
	})
	return cronInstance
}
