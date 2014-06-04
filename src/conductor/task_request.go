// task_request.go
//
package main

import (
	"sort"
	"time"
	o "orchestra"
)

type TaskRequest struct {
	job		*JobRequest
	Player		string				`json:"player"`
	State		TaskState			`json:"state"`
	RetryTime	time.Time			`json:"retrytime"`
}

func NewTaskRequest() (tr *TaskRequest) {
	tr = new(TaskRequest)
	tr.State = TASK_QUEUED

	return tr
}

func (task *TaskRequest) Encode() (ptr *o.TaskRequest) {
	ptr = new(o.TaskRequest)
	ptr.JobName = task.job.Score
	ptr.Id = task.job.Id
	ptr.Parameters = o.CopyMap(task.job.Params)

	return ptr
}

func (task *TaskRequest) IsTarget(player string) (valid bool) {
	valid = false
	if task.Player == "" {
		n := sort.SearchStrings(task.job.Players, player)
		if task.job.Players[n] == player {
			valid = true
		}
	} else {
		if task.Player == player {
			valid = true
		}
	}
	return valid
}
