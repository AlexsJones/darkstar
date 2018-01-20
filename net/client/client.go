package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"
)

//Configuration ...
type Configuration struct {
	Address   string
	Port      int
	CertPath  string
	KeyPath   string
	Message   string
	SleepTime time.Duration
}

//Run ...
func Run(clientconfig *Configuration) {
	cert, err := tls.LoadX509KeyPair(clientconfig.CertPath, clientconfig.KeyPath)
	if err != nil {
		panic(err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}

	for {
		conn, err := tls.Dial("tcp", clientconfig.Address+":"+strconv.Itoa(clientconfig.Port), &config)
		if err != nil {
			panic(err)
		}
		log.Println("client: connected to: ", conn.RemoteAddr())

		state := conn.ConnectionState()
		for _, v := range state.PeerCertificates {
			fmt.Println(x509.MarshalPKIXPublicKey(v.PublicKey))
			fmt.Println(v.Subject)
		}
		log.Println("client: handshake: ", state.HandshakeComplete)
		log.Println("client: mutual: ", state.NegotiatedProtocolIsMutual)

		n, err := io.WriteString(conn, clientconfig.Message)
		if err != nil {
			log.Fatalf("client: write: %s", err)
		}
		log.Printf("client: wrote %q (%d bytes)", clientconfig.Message, n)

		reply := make([]byte, 2048)
		n, err = conn.Read(reply)
		if err != nil {
			log.Fatalf("client: write: %s", err)
		}
		log.Printf("Receieved reply of %d bytes\n", n)
		//Process reply
		ServerHandler(reply[:n], conn)

		conn.Close()

		time.Sleep(clientconfig.SleepTime)
	}

}
