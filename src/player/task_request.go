// task_request.go
//

package main

import (
	o "orchestra"
)

type TaskRequest struct {
	Id		uint64				`json:"id"`
	Score		string				`json:"score"`
	State		TaskState			`json:"state"`		
	Params		map[string]string		`json:"params"`
	MyResponse	*TaskResponse			`json:"response"`
	RetryTime	int64				`json:"retrytime"`
}

func NewTaskRequest() (req *TaskRequest) {
	req = new(TaskRequest)
	return req
}

/* Map a wire task to an internal Task Request.
*/
func TaskFromProto(ptr *o.ProtoTaskRequest) (t *TaskRequest) {
	t = NewTaskRequest()
	
	t.Score = *(ptr.Jobname)
	t.Id = *(ptr.Id)
	t.Params = o.MapFromProtoJobParameters(ptr.Parameters)

	return t
}

