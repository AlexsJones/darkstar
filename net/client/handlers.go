package client

import (
	"log"

	"github.com/AlexsJones/darkstar/net/data"
	"github.com/AlexsJones/darkstar/net/data/message"
)

//ServerHandler reads instruction sets ...
func ServerHandler(raw []byte) {

	iface, t := data.TryUnmarshal(raw)
	switch t {

	case data.ProtoMessage:
		log.Println("Received message")
		m := iface.(*message.Message)

		ins := m.GetCurrentInstruction()
		if ins != nil {

		} else {
			log.Println("Receieved message had no current instruction")
		}
	case data.ProtoUnknown:
		log.Println("Protocol received unknown")
	}
}
