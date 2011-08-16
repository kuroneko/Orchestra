// task_request.go
//
package main

import (
	"sort"
	o "orchestra"
)

type TaskRequest struct {
	job		*JobRequest
	Player		string				`json:"player"`
	State		TaskState			`json:"state"`
	RetryTime	int64				`json:"retrytime"`
}

func (task *TaskRequest) Encode() (ptr *o.ProtoTaskRequest) {
	ptr = new(o.ProtoTaskRequest)
	ptr.Jobname = &task.job.Score
	ptr.Id = new(uint64)
	*ptr.Id = task.job.Id
	ptr.Parameters = o.ProtoJobParametersFromMap(task.job.Params)

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
