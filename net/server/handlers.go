package server

import (
	"fmt"
	"log"
	"net"

	"github.com/AlexsJones/darkstar/database/actor"
	"github.com/AlexsJones/darkstar/net/data"
	"github.com/AlexsJones/darkstar/net/data/message"
	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	model "gopkg.in/jeevatkm/go-model.v1"
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
			message := iface.(*message.Message)
			var ac actor.Actor

			//Check for existing actor ----------------------------------------------
			if err := databaseConnection.Find(&ac, actor.Actor{ActorID: message.ActorID}).Error; err != nil {

				color.Red(err.Error())

				//Map actor -------------------------------------------------------------
				var newactor actor.Actor
				errs := model.Copy(&newactor, message)
				if errs != nil {
					log.Println(errs)
				}

				log.Printf("New actor has connected to darkstar %+v\n", newactor)
				//-----------------------------------------------------------------------
				//Insert actor ----------------------------------------------------------
				if err := databaseConnection.Create(&newactor).Error; err != nil {
					color.Red(err.Error())
				}
				//-----------------------------------------------------------------------
				ac = newactor
			} else {
				log.Printf("Actor has reconnected %+v\n", ac)
			}
			//Reply is based from the incoming message with modified sub object -------
			//This makes sure each pulse from the client has the lastest module set ---
			ins := data.UpgradeMessage(message, serverConfig.ModuleName)
			n, _ := conn.Write([]byte(ins))
			log.Printf("Wrote %d byte response\n", n)
			//-------------------------------------------------------------------------
			log.Println("server: conn: closed")
			// Analyze any payload data received --------------------------------------

			if message.CurrentInstruction != nil {
				log.Printf("Reading actor instruction %s payload\n", message.CurrentInstruction.ModuleName)
				if message.CurrentInstruction.ModulePayload != nil {
					color.Green(fmt.Sprintf("Actor %s module %s payload has been received\n", ac.ActorID, message.CurrentInstruction.ModuleName))
				} else {
					color.Red(fmt.Sprintf("Actor %s module %s payload was empty\n", ac.ActorID, message.CurrentInstruction.ModuleName))
				}

			}
			//-----------------------------------------------------------------------
		case data.ProtoUnknown:
			color.Red("Receieved an unknown message type")
		}
	}
}
