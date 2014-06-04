// task_request.go
//

package main

import (
	o "orchestra"
)

type TaskRequest struct {
	Id		uint64				`json:"id"`
	Score		string				`json:"score"`
	Params		map[string]string		`json:"params"`
	MyResponse	*TaskResponse			`json:"response"`
}

func NewTaskRequest() (req *TaskRequest) {
	req = new(TaskRequest)
	return req
}

/* Map a wire task to an internal Task Request.
*/
func TaskFromProto(ptr *o.TaskRequest) (t *TaskRequest) {
	t = NewTaskRequest()
	
	t.Score = ptr.JobName
	t.Id = ptr.Id
	t.Params = o.CopyMap(ptr.Parameters)

	return t
}

