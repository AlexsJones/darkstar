package server

import (
	"fmt"
	"log"
	"net"

	"github.com/AlexsJones/darkstar/database"
	"github.com/AlexsJones/darkstar/net/data"
	"github.com/AlexsJones/darkstar/net/data/message"
	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
)

//ClientHandler is the behaviour on initial request to server
func ClientHandler(databaseConnection *gorm.DB, conn net.Conn, serverConfig *Configuration) {
	defer conn.Close()
	buf := make([]byte, 512)
	for {
		log.Print("server: conn: waiting")
		n, err := conn.Read(buf)
		if err != nil {
			if err != nil {
				log.Printf("server: conn: read: %s", err)
			}
			break
		}
		//TODO move this into another module

		iface, t := data.TryUnmarshal(buf[:n])
		switch t {
		case data.ProtoMessage:
			msg := iface.(*message.Message)

			color.Yellow(fmt.Sprintf("%+v\n", msg))
			//Check for existing actor ----------------------------------------------
			if err := databaseConnection.Find(&msg, message.Message{ActorID: msg.ActorID}).Error; err != nil {

				color.Red(err.Error())

				log.Printf("New actor has connected to darkstar %+v\n", msg)

				//-----------------------------------------------------------------------
				//Insert actor ----------------------------------------------------------
				if err := databaseConnection.Create(&msg).Error; err != nil {
					color.Red(err.Error())
				}
				//-----------------------------------------------------------------------

			} else {
				log.Printf("Actor has reconnected %+v\n", msg)

				//Check reconnected actor for module data

				if msg.GetCurrentInstruction() != nil {
					if len(msg.GetCurrentInstruction().ModulePayload.Data) > 0 {
						color.Green(fmt.Sprintf("Retrieved Actor %s module %s payload data\n", msg.ActorID, msg.CurrentInstruction.ModuleName))

						if err := databaseConnection.Create(&database.Module{ActorID: msg.ActorID, ModuleName: msg.CurrentInstruction.ModuleName, Data: msg.CurrentInstruction.ModulePayload.Data}).Error; err != nil {
							color.Red(err.Error())
						}
					}
				}
			}
			//-------------------------------------------------------------------------
			//Reply is based from the incoming message with modified sub object -------
			//This makes sure each pulse from the client has the lastest module set ---
			ins := data.MessageAddInstruction(msg, serverConfig.ModuleName)
			n, _ := conn.Write([]byte(ins))
			log.Printf("Wrote %d byte response\n", n)
			//-------------------------------------------------------------------------
			log.Println("server: conn: closed")
			// Analyze any payload data received --------------------------------------

			//-------------------------------------------------------------------------
		case data.ProtoUnknown:
			color.Red("Receieved an unknown message type")
		}
	}
}
