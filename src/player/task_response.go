// task_response.go
//
package main

import (
	"time"
	o "orchestra"
)

type TaskResponse struct {
	id		uint64				
	State		ResponseState			`json:"state"`
	Response	map[string]string		`json:"response"`
	// player only fields
	RetryTime	time.Time			`json:"retrytime"`
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

func (resp *TaskResponse) Encode() (ptr *o.TaskResponse) {
	ptr = new(o.TaskResponse)
	
	switch resp.State {
	case RESP_RUNNING:
		ptr.Status = o.TaskStatus_InProgress
	case RESP_FINISHED:
		ptr.Status = o.TaskStatus_Success
	case RESP_FAILED:
		ptr.Status = o.TaskStatus_Failed
	case RESP_FAILED_UNKNOWN_SCORE:
		ptr.Status = o.TaskStatus_Unknown
	case RESP_FAILED_HOST_ERROR:
		ptr.Status = o.TaskStatus_HostFailure
	case RESP_FAILED_UNKNOWN:
		ptr.Status = o.TaskStatus_UnknownFailure
	}
	ptr.Id = resp.id
	ptr.Response = o.CopyMap(resp.Response)

	return ptr
}

