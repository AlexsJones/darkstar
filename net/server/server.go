package server

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/AlexsJones/darkstar/net/data/instruction"
	"github.com/gogo/protobuf/proto"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

//Configuration ...
type Configuration struct {
	Address       string
	Port          int
	CertPath      string
	KeyPath       string
	ClientHandler func(Database *gorm.DB, conn net.Conn, serverConfig *Configuration)
	Mode          string
	Database      *gorm.DB
}

//Start ...
func Start(serverConfig *Configuration) error {

	//Load certs -----------------------------------------------------------------
	cert, err := tls.LoadX509KeyPair(serverConfig.CertPath, serverConfig.KeyPath)
	if err != nil {
		return err
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}}
	// ---------------------------------------------------------------------------
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
		go serverConfig.ClientHandler(serverConfig.Database, conn, serverConfig)
	}
	return nil
}

//CreateInstruction ....
func CreateInstruction(serverConfig *Configuration) string {

	ins := &instruction.Instruction{ModuleName: serverConfig.Mode}
	out, err := proto.Marshal(ins)
	if err != nil {
		log.Printf(err.Error())
		os.Exit(1)
	}
	return string(out)

}
