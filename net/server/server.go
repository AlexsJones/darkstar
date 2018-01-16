package server

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"strconv"
)

//Configuration ...
type Configuration struct {
	Address       string
	Port          int
	CertPath      string
	KeyPath       string
	ClientHandler func(conn net.Conn)
}

//Start ...
func Start(serverConfig *Configuration) error {
	cert, err := tls.LoadX509KeyPair(serverConfig.CertPath, serverConfig.KeyPath)
	if err != nil {
		return err
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}}
	config.Rand = rand.Reader
	service := serverConfig.Address + ":" + strconv.Itoa(serverConfig.Port)
	listener, err := tls.Listen("tcp", service, &config)
	if err != nil {
		return err
	}
	log.Print("server: listening")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)
			break
		}
		defer conn.Close()
		log.Printf("server: accepted from %s", conn.RemoteAddr())
		tlscon, ok := conn.(*tls.Conn)
		if ok {
			log.Print("ok=true")
			state := tlscon.ConnectionState()
			for _, v := range state.PeerCertificates {
				log.Print(x509.MarshalPKIXPublicKey(v.PublicKey))
			}
		}
		go serverConfig.ClientHandler(conn)
	}
	return nil
}
