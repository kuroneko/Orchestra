// resp_state.go

package main

import (
	"os"
	"json"
)

type ResponseState int

const (
	// Response states
	RESP_PENDING			= ResponseState(iota)	// internal state, not wire.
	RESP_RUNNING
	RESP_FINISHED
	RESP_FAILED
	RESP_FAILED_UNKNOWN_SCORE
	RESP_FAILED_HOST_ERROR
	RESP_FAILED_UNKNOWN
)

func (rs ResponseState) String() (strout string) {
	switch rs {
	case RESP_RUNNING:
		return "PENDING"
	case RESP_FINISHED:
		return "OK"
	case RESP_FAILED:
		return "FAIL"
	case RESP_FAILED_UNKNOWN_SCORE:
		return "UNK_SCORE"
	case RESP_FAILED_HOST_ERROR:
		return "HOST_ERROR"
	case RESP_FAILED_UNKNOWN:
		return "UNKNOWN_FAILURE"
	}
	return ""
}

func (rs ResponseState) MarshalJSON() (out []byte, err os.Error) {
	strout := rs.String()
	if strout != "" {
		return json.Marshal(strout)
	}
	return nil, InvalidValueError
}

func (rs ResponseState) UnmarshalJSON(in []byte) (err os.Error) {
	var statestr string
	err = json.Unmarshal(in, &statestr)
	if err != nil {
		return err
	}
	switch statestr {
	case "PENDING":
		rs = RESP_PENDING
	case "OK":
		rs = RESP_FINISHED
	case "FAIL":
		rs = RESP_FAILED
	case "UNK_SCORE":
		rs = RESP_FAILED_UNKNOWN_SCORE
	case "HOST_ERROR":
		rs = RESP_FAILED_HOST_ERROR
	case "UNKNOWN_FAILURE":
		rs = RESP_FAILED_UNKNOWN
	default:
		return InvalidValueError
	}
	return nil
}

