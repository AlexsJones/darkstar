package server

import (
	"log"
	"net"

	"github.com/AlexsJones/darkstar/data/message"
	"github.com/gogo/protobuf/proto"
)

//ClientHandler is the behaviour on initial request to server
func ClientHandler(conn net.Conn) {
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
		log.Println(message)
	}
	log.Println("server: conn: closed")

}
