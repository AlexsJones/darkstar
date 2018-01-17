package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/AlexsJones/darkstar/net"
	"github.com/AlexsJones/darkstar/net/data/message"
	"github.com/gogo/protobuf/proto"
	"github.com/matishsiao/goInfo"
	uuid "github.com/nu7hatch/gouuid"
)

//Configuration ...
type Configuration struct {
	Address  string
	Port     int
	CertPath string
	KeyPath  string
	Message  string
}

//Send ...
func Send(clientconfig *Configuration) error {
	cert, err := tls.LoadX509KeyPair(clientconfig.CertPath, clientconfig.KeyPath)
	if err != nil {
		return err
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", clientconfig.Address+":"+strconv.Itoa(clientconfig.Port), &config)
	if err != nil {
		return err
	}
	defer conn.Close()
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

	reply := make([]byte, 256)
	n, err = conn.Read(reply)
	if err != nil {
		log.Fatalf("client: write: %s", err)
	}
	log.Printf("client: read %q (%d bytes)", string(reply[:n]), n)
	log.Print("client: exiting")

	return nil
}

//CreateMessage ...
func CreateMessage() string {

	gi := goInfo.GetInfo()
	message := &message.Message{}
	u, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}

	message.ActorID = u.String()
	message.IPAddress = net.GetLocalIP().String()
	message.GoOS = gi.GoOS
	message.Core = gi.Core
	message.Hostname = gi.Hostname
	message.Kernel = gi.Kernel
	message.OS = gi.OS
	message.Platform = gi.Platform
	message.CPUs = int32(gi.CPUs)

	out, err := proto.Marshal(message)
	if err != nil {
		log.Printf(err.Error())
		os.Exit(1)
	}
	return string(out)
}
