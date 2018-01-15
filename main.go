package main

import (
	"flag"
	"log"
	"net"

	"github.com/AlexsJones/darkstar/client"
	"github.com/AlexsJones/darkstar/server"
)

func main() {
	var mode = flag.String("mode", "server", "Darkstar run mode")
	var clientkeypath = flag.String("clientkeypath", "certs/client.key", "Client key path")
	var clientcertpath = flag.String("clientcertpath", "certs/client.pem", "Client cert path")
	var clientPort = flag.Int("clientport", 8080, "Client port")
	var serverkeypath = flag.String("serverkeypath", "certs/server.key", "Server key path")
	var servercertpath = flag.String("servercertpath", "certs/server.pem", "Server cert path")
	var serverPort = flag.Int("serverport", 8080, "Server port")
	flag.Parse()

	switch *mode {
	case "client":

		config := &client.Configuration{Message: "This is a test",
			Address: "0.0.0.0", CertPath: *clientcertpath, KeyPath: *clientkeypath, Port: *clientPort}
		client.Send(config)
	default:
		conf := &server.Configuration{Address: "0.0.0.0", CertPath: *servercertpath, KeyPath: *serverkeypath,
			Port: *serverPort,
			ClientHandler: func(conn net.Conn) {
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
					log.Printf("server: conn: echo %q\n", string(buf[:n]))
					n, err = conn.Write(buf[:n])

					log.Printf("server: conn: wrote %d bytes", n)

					if err != nil {
						log.Printf("server: write: %s", err)
						break
					}
				}
				log.Println("server: conn: closed")

			}}

		if err := server.Start(conf); err != nil {
			log.Printf(err.Error())
		}
	}
}
