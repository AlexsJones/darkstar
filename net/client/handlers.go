package client

import (
	"log"

	"github.com/AlexsJones/darkstar/net/data"
	"github.com/AlexsJones/darkstar/net/data/instruction"
)

//ServerHandler reads instruction sets ...
func ServerHandler(raw []byte) {

	iface, t := data.TryUnmarshal(raw)
	switch t {
	case data.ProtoInstruction:
		log.Println("Received instruction")
		instruction := iface.(*instruction.Instruction)
		log.Println(instruction)
	}
}
