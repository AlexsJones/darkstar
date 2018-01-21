package data

import (
	"log"
	"os"

	"github.com/AlexsJones/darkstar/net"
	"github.com/AlexsJones/darkstar/net/data/message"
	"github.com/gogo/protobuf/proto"
	"github.com/matishsiao/goInfo"
	uuid "github.com/nu7hatch/gouuid"
)

//Protocoltype ...
type Protocoltype int

const (
	ProtoMessage Protocoltype = 0
	ProtoUnknown Protocoltype = 3
)

//TryUnmarshal attempts to unmarshal code
func TryUnmarshal(raw []byte) (interface{}, Protocoltype) {

	m := &message.Message{}
	if err := proto.Unmarshal(raw, m); err == nil {
		log.Println("Unmarshalled Message")
		return m, ProtoMessage
	}
	return nil, ProtoUnknown
}

//CreateMessage ...
func CreateMessage() string {

	gi := goInfo.GetInfo()
	message := &message.Message{}
	u, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}

	message.ActorID = u.String()
	message.IPAddress = net.GetLocalIP().String()
	message.GoOS = gi.GoOS
	message.Core = gi.Core
	message.Hostname = gi.Hostname
	message.Kernel = gi.Kernel
	message.OS = gi.OS
	message.Platform = gi.Platform
	message.CPUs = int32(gi.CPUs)

	out, err := proto.Marshal(message)
	if err != nil {
		log.Printf(err.Error())
		os.Exit(1)
	}
	return string(out)
}

//MessageAddInstruction adds additional information into a message
func MessageAddInstruction(msg *message.Message, moduleName string) string {

	//Add the instruction
	msg.CurrentInstruction = &message.Message_Instruction{ModuleName: moduleName}
	out, err := proto.Marshal(msg)
	if err != nil {
		log.Printf(err.Error())
		os.Exit(1)
	}
	return string(out)
}

//MessageAddPayload adds additional information into a message
func MessageAddPayload(msg *message.Message, payload string) string {

	//Add the instruction
	msg.CurrentInstruction.ModulePayload = &message.Message_Instruction_Payload{Data: payload}
	out, err := proto.Marshal(msg)
	if err != nil {
		log.Printf(err.Error())
		os.Exit(1)
	}
	return string(out)
}
