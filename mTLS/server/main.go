package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sync"
)

const (
	PORT = ":55555" //may need this in shared folder with client
	user = "user"
	pass = "user"

	CLIENT_PUBLIC_CERT = "../keys/client.signed.cert.pem"
	SERVER_PUBLIC_KEY  = "../keys/server.signed.cert.pem"
	SERVER_PRIVATE_KEY = "../keys/server.pem"
	CA_CHAIN           = "../keys/ca-chain.cert.pem"
)

var (
	validPath = regexp.MustCompile("^/command/(\\d+)?$")

	// Ideally this would be in a database rather than memory
	runLog = make(map[int]*Job)

	counter = SafeCounter{id: 0}
)

type SafeCounter struct {
	mu sync.Mutex
	id int
}

func (c *SafeCounter) Increment() int {
	c.mu.Lock()
	c.id++
	newId := c.id
	c.mu.Unlock()

	return newId
}

func addMTLS(server *http.Server) {
	//using self signed certs
	clientCert, err := ioutil.ReadFile(CA_CHAIN) //CLIENT_PUBLIC_CERT)
	if err != nil {
		log.Fatal(err)
	}

	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM(clientCert)

	tlsConfig := &tls.Config{
		ClientAuth:               tls.RequireAndVerifyClientCert,
		ClientCAs:                clientCertPool,
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS13,
	}

	tlsConfig.BuildNameToCertificate()

	server.TLSConfig = tlsConfig
}

func main() {
	server := &http.Server{Addr: PORT}
	addMTLS(server)

	http.HandleFunc("/command/", routeHandler)
	//log.Fatal(http.ListenAndServe(PORT, nil))

	err := server.ListenAndServeTLS(SERVER_PUBLIC_KEY, SERVER_PRIVATE_KEY)
	if err != nil {
		log.Fatal(err)
	}
}
