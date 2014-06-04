// task_response.go
//
package main

import (
	o "orchestra"
)

type TaskResponse struct {
	id		uint64				
	State		ResponseState			`json:"state"`
	Response	map[string]string		`json:"response"`	
}

// Response related magic

func NewTaskResponse() (resp *TaskResponse) {
	resp = new(TaskResponse)
	resp.Response = make(map[string]string)

	return resp
}

func (resp *TaskResponse) IsFinished() bool {
	return resp.State.Finished()
}

func (resp *TaskResponse) DidFail() bool {
	return resp.State.Failed()
}

func (resp *TaskResponse) CanRetry() bool {
	return resp.State.CanRetry()
}


func ResponseFromProto(ptr *o.TaskResponse) (r *TaskResponse) {
	r = new(TaskResponse)

	switch (ptr.Status) {
	case o.TaskStatus_InProgress:
		r.State = RESP_RUNNING
	case o.TaskStatus_Success:
		r.State = RESP_FINISHED
	case o.TaskStatus_Failed:
		r.State = RESP_FAILED
	case o.TaskStatus_HostFailure:
		r.State = RESP_FAILED_HOST_ERROR
	case o.TaskStatus_Unknown:
		r.State = RESP_FAILED_UNKNOWN_SCORE
	case o.TaskStatus_UnknownFailure:
		fallthrough
	default:
		r.State = RESP_FAILED_UNKNOWN
	}

	r.id = ptr.Id
	r.Response = o.CopyMap(ptr.Response)

	return r
}
