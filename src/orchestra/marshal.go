/* marshal.go
 *
 * Common wire marshalling code.
 */

package orchestra

import (
	"labix.org/v2/mgo/bson"
	"errors"
)

var (
	ErrUnknownType    = errors.New("Unknown Type in Encode request")
	ErrObjectTooLarge = errors.New("Encoded Object exceeds maximum encoding size")
)

const (
	MaximumPayloadSize = 0x10000
)

func (p *WirePkt) Decode() (obj interface{}, err error) {
	switch p.Type {
	case TypeNop:
		if p.Length != 0 {
			/* throw error later... */
			return nil, ErrMalformedMessage
		}
		return nil, nil
	case TypeIdentifyClient:
		ic := new(IdentifyClient)
		err := bson.Unmarshal(p.Payload[0:p.Length], ic)
		if err != nil {
			return nil, err
		}
		return ic, nil
	case TypeReadyForTask:
		if p.Length != 0 {
			/* throw error later... */
			return nil, ErrMalformedMessage
		}
		return nil, nil
	case TypeTaskRequest:
		tr := new(TaskRequest)
		err := bson.Unmarshal(p.Payload[0:p.Length], tr)
		if err != nil {
			return nil, err
		}
		return tr, nil
	case TypeTaskResponse:
		tr := new(TaskResponse)
		err := bson.Unmarshal(p.Payload[0:p.Length], tr)
		if err != nil {
			return nil, err
		}
		return tr, nil
	case TypeAcknowledgement:
		tr := new(Acknowledgement)
		err := bson.Unmarshal(p.Payload[0:p.Length], tr)
		if err != nil {
			return nil, err
		}
		return tr, nil
	}
	return nil, ErrUnknownMessage
}

func Encode(obj interface{}) (p *WirePkt, err error) {
	p = new(WirePkt)
	switch obj.(type) {
	case *IdentifyClient:
		p.Type = TypeIdentifyClient
	case *TaskRequest:
		p.Type = TypeTaskRequest
	case *TaskResponse:
		p.Type = TypeTaskResponse
	case *Acknowledgement:
		p.Type = TypeAcknowledgement
	default:
		Warn("Encoding unknown type!")
		return nil, ErrUnknownType
	}
	p.Payload, err = bson.Marshal(obj)
	if err != nil {
		return nil, err
	}
	if len(p.Payload) >= MaximumPayloadSize {
		return nil, ErrObjectTooLarge
	}
	p.Length = uint16(len(p.Payload))

	return p, nil
}

func MakeNop() (p *WirePkt) {
	p = new(WirePkt)
	p.Length = 0
	p.Type = TypeNop
	p.Payload = nil

	return p
}

func MakeIdentifyClient(hostname, clientid string) (p *WirePkt) {
	s := new(IdentifyClient)
	s.Hostname = hostname
	s.ClientId = clientid

	p, _ = Encode(s)

	return p
}

func MakeReadyForTask() (p *WirePkt) {
	p = new(WirePkt)
	p.Type = TypeReadyForTask
	p.Length = 0
	p.Payload = nil

	return p
}

/* We use the failure code for negative acknowledgements */
func MakeNack(id uint64) (p *WirePkt) {
	a := new(Acknowledgement)
	a.Id = id
	a.Response = AckType_Error
	p, _ = Encode(a)

	return p
}

// Construct a positive ACK for transmission
func MakeAck(id uint64) (p *WirePkt) {
	a := new(Acknowledgement)
	a.Id = id
	a.Response = AckType_OK
	p, _ = Encode(a)

	return p
}
