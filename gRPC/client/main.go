package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"time"

	"github.com/woodsjc/worker_api/gRPC/internal/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	USER    = "user"
	PASS    = "user"
	URL     = "https://localhost:55555/command"
	ADDRESS = "localhost:55555"

	CLIENT_PRIVATE_KEY = "../keys/client.pem"
	CLIENT_PUBLIC_KEY  = "../keys/client.signed.cert.pem"
	SERVER_PUBLIC_KEY  = "../keys/server.signed.cert.pem"
	CA_CHAIN           = "../keys/ca-chain.cert.pem"
)

func GetClientTransport() credentials.TransportCredentials {
	cert, err := tls.LoadX509KeyPair(CLIENT_PUBLIC_KEY, CLIENT_PRIVATE_KEY)
	if err != nil {
		log.Fatal(err)
	}

	caCert, err := ioutil.ReadFile(CA_CHAIN)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	transportCreds := credentials.NewTLS(&tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
	})

	return transportCreds
}

func main() {
	transportCreds := GetClientTransport()
	dialOption := grpc.WithTransportCredentials(transportCreds)
	conn, err := grpc.Dial(ADDRESS, dialOption)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	cli := worker.NewStartServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := cli.Start(ctx, &worker.StartRequest{Name: "ls", Args: "-ax"})
	if err != nil {
		log.Print(err)
	} else {
		log.Printf("Success resposne: %v", response)
	}

	response, err = cli.Start(ctx, &worker.StartRequest{Name: "env", Args: ""})
	if err != nil {
		log.Print(err)
	} else {
		log.Printf("Success resposne: %v", response)
	}
}
