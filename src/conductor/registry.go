// registry.go
//
// The Registry provides a 'threadsafe' interface to various global information stores.
//
// The registry dispatch thread is forbidden from performing any work that is likely to block.
// Result channels must be buffered with enough space for the full set of results.

package main

import (
	o "orchestra"
	"sort"
)

// Request Types
const (
	requestAddClient = iota
	requestGetClient
	requestDeleteClient
	requestSyncClients

	requestAddJob
	requestGetJob
	requestAddJobResult
	requestGetJobResult
	requestGetJobResultNames
	requestDisqualifyPlayer
	requestReviewJobStatus

	requestWriteJobUpdate
	requestWriteJobAll

	requestQueueSize		= 10
)

type registryRequest struct {
	operation		int
	id			uint64
	hostname		string
	hostlist		[]string
	job			*JobRequest
	tresp			*TaskResponse
	responseChannel		chan *registryResponse
}

type registryResponse struct {
	success			bool
	info			*ClientInfo
	tresp			*TaskResponse
	names			[]string
	jobs			[]*JobRequest
}

var (
	chanRegistryRequest	= make(chan *registryRequest, requestQueueSize)
 	clientList 		= make(map[string]*ClientInfo)
	jobRegister 		= make(map[uint64]*JobRequest)
)

func regInternalAdd(hostname string) {
	o.Warn("Registry: New Host \"%s\"", hostname)
	clientList[hostname] = NewClientInfo()
	// do this initialisation here since it'll help unmask sequencing errors
	clientList[hostname].pendingTasks = make(map[uint64]*TaskRequest)
	clientList[hostname].Player = hostname
}

func regInternalDel(hostname string) {
	o.Warn("Registry: Deleting Host \"%s\"", hostname)
	/* remove it from the registry */
	clientList[hostname] = nil, false
}

func manageRegistry() {
	for {
		req := <-chanRegistryRequest
		resp := new(registryResponse)
		/* by default, we failed. */
		resp.success = false
		switch (req.operation) {
		case requestAddClient:
			_, exists := clientList[req.hostname]
			if exists {
				resp.success = false
			} else {
				regInternalAdd(req.hostname)
				resp.success = true
			}
		case requestGetClient:
			clinfo, exists := clientList[req.hostname]
			if exists {
				resp.success = true
				resp.info = clinfo
			} else {
				resp.success = false
			}
		case requestDeleteClient:
			_, exists := clientList[req.hostname]
			if exists {
				resp.success = true
				regInternalDel(req.hostname)
			} else {
				resp.success = false
			}
		case requestSyncClients:
			/* we need to make sure the registered clients matches
			 * the hostlist we're given.
			 *
			 * First, we transform the array into a map
			 */
			newhosts := make(map[string]bool)
			for k,_ := range req.hostlist {
				newhosts[req.hostlist[k]] = true
			}
			/* now, scan the current list, checking to see if
			 * they exist.  Remove them from the newhosts map
			 * if they do exist. 
			 */
			for k,_ := range clientList {
				_, exists := newhosts[k]
				if exists {
					/* remove it from the newhosts map */
					newhosts[k] = false, false
				} else {
					regInternalDel(k)
				}
			}
			/* now that we're finished, we should only have
			 * new clients in the newhosts list left. 
			 */
			for k,_ := range newhosts {
				regInternalAdd(k)
			}
			/* and we're done. */
		case requestAddJob:
			if nil != req.job {
				// ensure that the players are sorted!
				sort.Strings(req.job.Players)
				// update the state
				req.job.updateState()
				// and register the job
				jobRegister[req.job.Id] = req.job
				// force a queue update.
				req.job.UpdateJobInformation()
				resp.success = true
			} else {
				resp.success = false
			}
		case requestGetJob:
			job, exists := jobRegister[req.id]
			resp.success = exists
			if exists {
				resp.jobs = make([]*JobRequest, 1)
				resp.jobs[0] = job
			}
		case requestAddJobResult:
			job, exists := jobRegister[req.tresp.id]
			resp.success = exists
			if exists {
				job.Results[req.hostname] = req.tresp
				// force a queue update.
				job.UpdateJobInformation()
			}
		case requestGetJobResult:
			job, exists := jobRegister[req.id]
			if exists {
				result, exists := job.Results[req.hostname]
				resp.success = exists
				if exists {
					resp.tresp = result
				}
			} else {
				resp.success = false
			}
		case requestGetJobResultNames:
			job, exists := jobRegister[req.id]
			resp.success = exists
			if exists {
				resp.names = make([]string, len(job.Results))
				idx := 0
				for k, _ := range job.Results {
					resp.names[idx] = k
					idx++
				}
			}
		case requestDisqualifyPlayer:
			job, exists := jobRegister[req.id]
			if exists {
				idx := sort.Search(len(job.Players), func(idx int) bool { return job.Players[idx] >= req.hostname })
				if (job.Players[idx] == req.hostname) {
					resp.success = true
					newplayers := make([]string, len(job.Players)-1)
					copy(newplayers[0:idx], job.Players[0:idx])
					copy(newplayers[idx:len(job.Players)-1], job.Players[idx+1:len(job.Players)])
					job.Players = newplayers
					job.updateState()
					// force a queue update.
					job.UpdateJobInformation()
				} else {
					resp.success = false
				}
			} else {
				resp.success = false
			}
		case requestReviewJobStatus:
			job, exists := jobRegister[req.id]
			resp.success = exists
			if exists {
				job.updateState()
				// force a queue update.
				job.UpdateJobInformation()
			}
		case requestWriteJobUpdate:
			job, exists := jobRegister[req.id]
			resp.success = exists
			if exists {
				job.UpdateJobInformation()
			}
		case requestWriteJobAll:
			for _, job := range jobRegister {
				job.UpdateJobInformation()
			}
			resp.success = true
		}
		if req.responseChannel != nil {
			req.responseChannel <- resp
		}
	}
}

func StartRegistry() {
	go manageRegistry()
}

func newRequest(wants_response bool) (req *registryRequest) {
	req = new(registryRequest)
	if wants_response {
		req.responseChannel = make(chan *registryResponse, 1)
	}

	return req
}
	
func ClientAdd(hostname string) (success bool) {
	r := newRequest(true)
	r.operation = requestAddClient
	r.hostname = hostname
	chanRegistryRequest <- r
	resp := <- r.responseChannel
	
	return resp.success
}

func ClientDelete(hostname string) (success bool) {
	r := newRequest(true)
	r.operation = requestDeleteClient
	r.hostname = hostname
	chanRegistryRequest <- r
	resp := <- r.responseChannel
	
	return resp.success
}

func ClientGet(hostname string) (info *ClientInfo) {
	r := newRequest(true)
	r.operation = requestGetClient
	r.hostname = hostname
	chanRegistryRequest <- r
	resp := <- r.responseChannel
	if resp.success {
		return resp.info
	}
	return nil
}

func ClientUpdateKnown(hostnames []string) {
	/* this is an asynchronous, we feed it into the registry 
	 * and it'll look after itself.
	*/
	r := newRequest(false)
	r.operation = requestSyncClients
	r.hostlist = hostnames
	chanRegistryRequest <- r
}

// Add a Job to the registry.  Return true if successful, returns
// false if the job is lacking critical information (such as a JobId)
// and can't be registered.
func JobAdd(job *JobRequest) bool {
	rr := newRequest(true)
	rr.operation = requestAddJob
	rr.job = job

	chanRegistryRequest <- rr
	resp := <- rr.responseChannel 
	return resp.success
}

// Get a Job from the registry.  Returns the job if successful,
// returns nil if the job couldn't be found.
func JobGet(id uint64) *JobRequest {
	rr := newRequest(true)
	rr.operation = requestGetJob
	rr.id = id

	chanRegistryRequest <- rr
	resp := <- rr.responseChannel
	if resp.jobs == nil {
		return nil
	}
	return resp.jobs[0]
}

// Attach a result to a Job in the Registry
//
// This exists in order to prevent nasty concurrency problems
// when trying to put results back onto the job.  Reading a job is far
// less of a problem than writing to it.
func JobAddResult(playername string, task *TaskResponse) bool {
	rr := newRequest(true)
	rr.operation = requestAddJobResult
	rr.tresp = task
	rr.hostname = playername
	chanRegistryRequest <- rr
	resp := <- rr.responseChannel
	return resp.success
}

// Get a result from the registry
func JobGetResult(id uint64, playername string) (tresp *TaskResponse) {
	rr := newRequest(true)
	rr.operation = requestGetJobResult
	rr.id = id
	rr.hostname = playername
	chanRegistryRequest <- rr
	resp := <- rr.responseChannel
	return resp.tresp
}

// Get a list of names we have results for against a given job.
func JobGetResultNames(id uint64) (names []string) {
	rr := newRequest(true)
	rr.operation = requestGetJobResultNames
	rr.id = id

	chanRegistryRequest <- rr
	resp := <- rr.responseChannel 
	return resp.names
}

//  Disqualify a player from servicing a job
func JobDisqualifyPlayer(id uint64, playername string) bool {
	rr := newRequest(true)
	rr.operation = requestDisqualifyPlayer
	rr.id = id
	rr.hostname = playername

	chanRegistryRequest <- rr
	resp := <- rr.responseChannel

	return resp.success
}

func JobReviewState(id uint64) bool {
	rr := newRequest(true)
	rr.operation = requestReviewJobStatus
	rr.id = id

	chanRegistryRequest <- rr
	resp := <- rr.responseChannel

	return resp.success
}

func JobWriteUpdate(id uint64) {
	rr := newRequest(false)
	rr.operation = requestWriteJobUpdate
	rr.id = id
	chanRegistryRequest <- rr
}

func JobWriteAll() bool {
	rr := newRequest(true)
	rr.operation = requestWriteJobAll

	chanRegistryRequest <- rr
	resp := <-rr.responseChannel

	return resp.success
}

// Ugh.
func (job *JobRequest) updateState() {
	if job.Results == nil {
		o.Assert("job.Results nil for jobid %d", job.Id)
		return
	}
	switch job.Scope {
	case SCOPE_ONEOF:
		// look for a success (any success) in the responses
		var success bool = false
		for host, res := range job.Results {
			if res == nil {
				o.Debug("nil result for %s?", host)
				continue
			}
			if res.State == RESP_FINISHED {
				success = true
				break
			}
		}
		// update the job state based upon these findings
		if success {
			job.State = JOB_SUCCESSFUL
		} else {
			if len(job.Players) < 1 {
				job.State = JOB_FAILED
			} else {
				job.State = JOB_PENDING
			}
		}
	case SCOPE_ALLOF:
		var success int = 0
		var failed  int = 0
		
		for pidx := range job.Players {
			p := job.Players[pidx]
			resp, exists := job.Results[p]
			if exists {
				if resp.DidFail() {
					failed++
				} else if resp.State == RESP_FINISHED {
					success++
				}
			}
		}
		if (success + failed) < len(job.Players) {
			job.State = JOB_PENDING
		} else if success == len(job.Players) {
			job.State = JOB_SUCCESSFUL
		} else if failed == len(job.Players) {
			job.State = JOB_FAILED
		} else {
			job.State = JOB_FAILED_PARTIAL
		}
	}
}
