package database

import "time"

type AutomationExecution struct {
	Id            uint64
	Automation_id uint64
	Success       bool
	Run_datetime  time.Time
}
