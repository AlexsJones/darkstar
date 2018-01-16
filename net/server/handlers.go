package server

import (
	"log"
	"net"

	"github.com/AlexsJones/darkstar/database/actor"
	"github.com/AlexsJones/darkstar/net/data/message"
	"github.com/gogo/protobuf/proto"
	"github.com/jinzhu/gorm"
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

		var actor actor.Actor
		databaseConnection.First(&actor, "identifier = ?", message.ActorID)
		if &actor != nil {
			log.Printf("New actor has connected to darkstar %s\n", message.ActorID)
		} else {
			log.Printf("Actor %s has reconnected\n", message.ActorID)
		}
	}
	log.Println("server: conn: closed")

}
