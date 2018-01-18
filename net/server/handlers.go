package server

import (
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
			} else {
				log.Printf("Actor has reconnected %+v\n", ac)
			}
			//Reply -------------------------------------------------------------------
			ins := CreateInstruction(serverConfig)

			conn.Write([]byte(ins))
			//-------------------------------------------------------------------------
			log.Println("server: conn: closed")

		case data.ProtoOperation:

		}
	}
}
