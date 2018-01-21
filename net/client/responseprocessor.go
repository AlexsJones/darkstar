package client

import (
	"fmt"
	"log"

	"github.com/AlexsJones/darkstar/modules"
	"github.com/AlexsJones/darkstar/net/data"
	"github.com/AlexsJones/darkstar/net/data/message"
	"github.com/fatih/color"
)

//ResponseProcessor reads instruction sets ...
func ResponseProcessor(raw []byte) (string, bool) {

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
			d := modules.IModule.Execute(iface)

			msg := data.MessageAddPayload(m, d)

			return msg, true

		} else {
			log.Println("Receieved message had no current instruction")
		}
	case data.ProtoUnknown:
		log.Println("Protocol received unknown")
	}

	return "", false
}
