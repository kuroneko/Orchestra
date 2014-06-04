/* msg.go
 *
 * BSON message formats
 */
package orchestra

/* P->C : Provide Client Identity and negotiate other initial parameters */
type IdentifyClient struct {
	// the client's hostname
	Hostname string		`bson:"hostname"`
	// optional player version information
	ClientId string		`bson:"version,omitempty"`
}

/* C->P : Do Shit kthxbye */
type TaskRequest struct {
	JobName string					`bson:"job_name"`
	Id uint64						`bson:"id"`
	Parameters map[string] string	`bson:"params"`
}

const (
	AckType_OK = uint8(1)
	AckType_Error = uint8(3)
)

/* C->P, P->C : Acknowledge Message */
type Acknowledgement struct {
	Id uint64			`bson:"id"`
	Response uint8		`bson:"type"`
}

const (
	TaskStatus_InProgress = uint8(2)
	TaskStatus_Success = uint8(3)
	TaskStatus_Failed = uint8(4)
	TaskStatus_HostFailure = uint8(5)
	TaskStatus_Unknown = uint8(6)
	TaskStatus_UnknownFailure = uint8(7)
)

/* P->C : Results from Task */
type TaskResponse struct {
	Id uint64					`bson:"id"`
	Status uint8				`bson:"status"`
	Response map[string] string	`bson:"response"`
}

func CopyMap(in map[string] string) (out map[string] string) {
	out = make(map[string] string)
	for k, v := range in {
		out[k] = v
	}
	return out
}