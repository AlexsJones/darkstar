package data

import (
	"github.com/AlexsJones/darkstar/net/data/message"
	"github.com/AlexsJones/darkstar/net/data/operation"
	"github.com/gogo/protobuf/proto"
)

//Protocoltype ...
type Protocoltype int

const (
	ProtoMessage Protocoltype = iota
	ProtoOperation
	ProtoUnknown
)

//TryUnmarshal attempts to unmarshal code
func TryUnmarshal(raw []byte) (interface{}, Protocoltype) {

	m := &message.Message{}
	if err := proto.Unmarshal(raw, m); err == nil {
		return m, ProtoMessage
	}
	o := &operation.Operation{}
	if err := proto.Unmarshal(raw, o); err == nil {
		return o, ProtoOperation
	}

	return nil, ProtoUnknown
}
