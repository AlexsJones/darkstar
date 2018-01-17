package server

import (
	"log"
	"net"

	"github.com/AlexsJones/darkstar/database/actor"
	"github.com/AlexsJones/darkstar/net/data/message"
	"github.com/gogo/protobuf/proto"
	"github.com/jinzhu/gorm"
	model "gopkg.in/jeevatkm/go-model.v1"
)

//ClientHandler is the behaviour on initial request to server
func ClientHandler(databaseConnection *gorm.DB, conn net.Conn) {
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
		message := &message.Message{}
		if err := proto.Unmarshal(buf[:n], message); err != nil {
			log.Printf(err.Error())
			return
		}

		var ac actor.Actor
		databaseConnection.First(&ac, "id = ?", message.ActorID)
		if &ac != nil {
			log.Printf("New actor has connected to darkstar %s:%s\n", message.ActorID, message.IPAddress)
			//Map actor -----------------------------------------------------------------------
			var newactor actor.Actor
			errs := model.Copy(&newactor, message)
			if errs != nil {
				log.Println(errs)
			}
			//---------------------------------------------------------------------------------
			//Insert actor --------------------------------------------------------------------
			databaseConnection.Create(&newactor)
			//---------------------------------------------------------------------------------
		} else {
			log.Printf("Actor %s:%s has reconnected\n", ac.ActorID, ac.IPAddress)
		}
	}
	log.Println("server: conn: closed")

}
