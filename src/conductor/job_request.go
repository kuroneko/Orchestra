// job_request.go
//

package main

import (
	"sort"
)

type JobRequest struct {
	Score		string				`json:"score"`
	Scope		JobScope			`json:"scope"`
	Players		[]string			`json:"players"`
	Id		uint64				`json:"id"`
	State		JobState			`json:"state"`		
	Params		map[string]string		`json:"params"`
	Tasks		[]*TaskRequest			`json:"tasks"`
	// These are private - you need to use the registry to access these
	results		map[string]*TaskResponse	`json:"results"`
}

func NewJobRequest() (req *JobRequest) {
	req = new(JobRequest)
	req.results = make(map[string]*TaskResponse)
	return req
}

func (req *JobRequest) normalise() {
	if (len(req.Players) > 1) {
		/* sort targets so search works */
		sort.Strings(req.Players)
	} else {
		if (req.Scope == SCOPE_ONEOF) {
			req.Scope = SCOPE_ALLOF
		}
	}
}

func (req *JobRequest) MakeTasks() (tasks []*TaskRequest) {
	req.normalise()

	var numtasks int
	
	switch (req.Scope) {
	case SCOPE_ONEOF:
		numtasks = 1
	case SCOPE_ALLOF:
		numtasks = len(req.Players)
	}
	tasks = make([]*TaskRequest, numtasks)
	
	for c := 0; c < numtasks; c++ {
		t := new(TaskRequest)
		t.State = TASK_QUEUED
		t.job = req
		if (req.Scope == SCOPE_ALLOF) {
			t.Player = req.Players[c]
		}
		tasks[c] = t
	}
	return tasks
}

func (req *JobRequest) Valid() bool {
	if (len(req.Players) <= 0) {
		return false
	}
	return true
}
