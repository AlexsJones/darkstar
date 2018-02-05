package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"golang.org/x/net/proxy"
)

//Configuration ...
type Configuration struct {
	Address string
	Port    int
	//TLS --------------------
	CertPath string
	KeyPath  string
	//SOCKS ------------------
	ProxyAddress string
	ProxyPort    int
	Auth         *proxy.Auth
	// -----------------------
	SleepTime time.Duration
}

//SendSOCKS ...
func SendSOCKS(clientconfig *Configuration, message string) ([]byte, int) {

	dialer, err := proxy.SOCKS5("tcp", clientconfig.ProxyAddress+":"+strconv.Itoa(clientconfig.ProxyPort), clientconfig.Auth, &net.Dialer{
		Timeout:   60 * time.Second,
		KeepAlive: 30 * time.Second,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
		os.Exit(1)
	}
	log.Printf("Dialling %s\n", clientconfig.Address+":"+strconv.Itoa(clientconfig.Port))
	conn, err := dialer.Dial("tcp", clientconfig.Address+":"+strconv.Itoa(clientconfig.Port))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	log.Println("client: connected to: ", conn.RemoteAddr())
	n, err := io.WriteString(conn, message)
	if err != nil {
		log.Fatalf("client: write: %s", err)
	}
	log.Printf("client: wrote %q (%d bytes)", message, n)

	reply := make([]byte, 2048)
	n, err = conn.Read(reply)
	if err != nil {
		log.Fatalf("client: write: %s", err)
	}
	log.Printf("Receieved reply of %d bytes\n", n)
	//Process reply
	return reply, n
}

//SendTLS ...
func SendTLS(clientconfig *Configuration, message string) ([]byte, int) {
	cert, err := tls.LoadX509KeyPair(clientconfig.CertPath, clientconfig.KeyPath)
	if err != nil {
		panic(err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}

	conn, err := tls.Dial("tcp", clientconfig.Address+":"+strconv.Itoa(clientconfig.Port), &config)
	if err != nil {
		panic(err)
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

	n, err := io.WriteString(conn, message)
	if err != nil {
		log.Fatalf("client: write: %s", err)
	}
	log.Printf("client: wrote %q (%d bytes)", message, n)

	reply := make([]byte, 2048)
	n, err = conn.Read(reply)
	if err != nil {
		log.Fatalf("client: write: %s", err)
	}
	log.Printf("Receieved reply of %d bytes\n", n)
	//Process reply
	return reply, n
}
