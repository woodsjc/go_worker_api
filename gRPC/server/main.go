package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"
	"regexp"
	"sync"

	"github.com/woodsjc/worker_api/gRPC/internal/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	validPath = regexp.MustCompile("^/command(/\\d+)?$")

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

type Server struct {
	worker.UnimplementedStartServiceServer
}

func (s *Server) Start(ctx context.Context,
	in *worker.StartRequest) (*worker.StartResponse, error) {
	log.Printf("Received: %v", in.GetName())

	id := counter.Increment()
	j := &Job{id: id}
	runLog[id] = j
	j.Start(in.Name, in.Args)

	return &worker.StartResponse{Id: int64(id)}, nil
}

func addMTLS() grpc.ServerOption {
	//using self signed certs
	serverCert, err := tls.LoadX509KeyPair(SERVER_PUBLIC_KEY, SERVER_PRIVATE_KEY)
	if err != nil {
		log.Fatal(err)
	}

	caChain, err := ioutil.ReadFile(CA_CHAIN) //CLIENT_PUBLIC_CERT)
	if err != nil {
		log.Fatal(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caChain)

	tlsConfig := &tls.Config{
		Certificates:             []tls.Certificate{serverCert},
		ClientAuth:               tls.RequireAndVerifyClientCert,
		ClientCAs:                caCertPool,
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS13,
	}

	tlsConfig.BuildNameToCertificate()
	return grpc.Creds(credentials.NewTLS(tlsConfig))
}

func main() {
	//server := &http.Server{Addr: PORT}
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatal(err)
	}
	serverOption := addMTLS()
	server := grpc.NewServer(serverOption)

	worker.RegisterStartServiceServer(server, &Server{})
	log.Fatal(server.Serve(listener))
}
