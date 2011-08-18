// job_request.go
//

package main

import (
	"sort"
	"json"
	"path"
	"os"
	"io"
	o "orchestra"
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
	Results		map[string]*TaskResponse	`json:"results"`
}

func NewJobRequest() (req *JobRequest) {
	req = new(JobRequest)
	req.Results = make(map[string]*TaskResponse)
	return req
}

func JobRequestFromReader(src io.Reader) (req *JobRequest, err os.Error) {
	req = NewJobRequest()
	jdec := json.NewDecoder(src)

	err = jdec.Decode(req)

	return req, err
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
		t := NewTaskRequest()
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

func (req *JobRequest) FilenameForSpool() string {
	if (req.State == JOB_PENDING) {
		return path.Join(GetSpoolDirectory(), "active", FilenameForJobId(req.Id))
	}
	return path.Join(GetSpoolDirectory(), "finished", FilenameForJobId(req.Id))
}

// dump the bytestream in buf into the serialisation file for req.
func (req *JobRequest) doSerialisation(buf []byte) {
	// first up, clean up old state.
	UnlinkNodesForJobId(req.Id)
	outpath := req.FilenameForSpool()
	fh, err := os.OpenFile(outpath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		o.Warn("Could not create persistence file %s: %s", outpath, err)
		return
	}
	defer fh.Close()
	fh.Write(buf)
}

func (req *JobRequest) UpdateJobInformation()  {
	buf, err := json.MarshalIndent(req, "", "  ")
	o.MightFail(err, "Failed to marshal job %d", req.Id)
	//FIXME: should try to do this out of the registry's thread.
	req.doSerialisation(buf)
}