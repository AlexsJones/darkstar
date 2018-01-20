package client

import (
	"fmt"
	"log"
	"net"

	"github.com/AlexsJones/darkstar/modules"
	"github.com/AlexsJones/darkstar/net/data"
	"github.com/AlexsJones/darkstar/net/data/message"
	"github.com/fatih/color"
)

//ServerHandler reads instruction sets ...
func ServerHandler(raw []byte, conn net.Conn) {

	iface, t := data.TryUnmarshal(raw)
	switch t {

	case data.ProtoMessage:
		log.Println("Received message")
		m := iface.(*message.Message)

		ins := m.GetCurrentInstruction()
		if ins != nil {

			//Load module to run
			iface, err := modules.LoadModule(ins.ModuleName)
			if err != nil {
				fmt.Errorf(err.Error())
			}
			color.Green(fmt.Sprintf("Loaded %s module\n", ins.ModuleName))
			//Execute module now has run
			data := modules.IModule.Execute(iface)

			conn.Write([]byte(data))

		} else {
			log.Println("Receieved message had no current instruction")
		}
	case data.ProtoUnknown:
		log.Println("Protocol received unknown")
	}
}
