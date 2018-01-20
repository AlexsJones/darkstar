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
		log.Println(m)
	case data.ProtoUnknown:
		log.Println("Protocol received unknown")
	}
}
