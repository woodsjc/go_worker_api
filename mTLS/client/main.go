package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	USER = "user"
	PASS = "user"
	URL  = "https://localhost:55555/command/"

	CLIENT_PRIVATE_KEY = "../keys/client.pem"
	CLIENT_PUBLIC_KEY  = "../keys/client.signed.cert.pem"
	SERVER_PUBLIC_KEY  = "../keys/server.signed.cert.pem"
	CA_CHAIN           = "../keys/ca-chain.cert.pem"
)

type Cli struct {
	*http.Client
}

type PostData struct {
	Name string
	Args string
}

type PostResponse struct {
	id int
}

type GetResponse struct {
	Id  int
	Log string
}

func setRequestHeaders(r *http.Request) {
	//r.SetBasicAuth(USER, PASS)
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Content-Type", "application/json")
}

func (c *Cli) postJob(job string, args string) {
	data, err := json.Marshal(PostData{job, args})
	if err != nil {
		log.Printf("Invalid json %v %v", job, args)
		return
	}

	request, err := http.NewRequest("POST", URL, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Unable to post job %v", err)
		return
	}

	setRequestHeaders(request)
	response, err := c.Do(request)
	if err != nil {
		log.Println(err)
	}

	if response.StatusCode == 200 {
		id := PostResponse{}
		log.Printf("Posted job: %v", json.NewDecoder(response.Body).Decode(&id))
	} else {
		log.Printf("Failed to post job %v %v: %v", job, args, response)
	}
}

func (c *Cli) getJob(id int) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s%d", URL, id), nil)
	if err != nil {
		log.Printf("Unable build new get request %v", err)
		return
	}

	setRequestHeaders(request)
	response, err := c.Do(request)
	if err != nil {
		log.Println(err)
	}

	if response.StatusCode == 200 {
		get := GetResponse{}
		err := json.NewDecoder(response.Body).Decode(&get)
		if err != nil {
			log.Printf("Unable to decode GET response: %v", get)
			return
		}
		log.Printf("Get job: %v", get)
	} else {
		log.Printf("Failed to get job %v: %v", id, response)
	}
}

func mTLSAuth() *Cli {
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

	cli := &Cli{&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{cert},
				MinVersion:   tls.VersionTLS13,
			},
		},
	}}
	return cli
}

func main() {
	cli := mTLSAuth()

	r, err := cli.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response from %v: %v", URL, r)

	cli.postJob("ls", "-ax")
	cli.postJob("env", "")

	cli.getJob(1)
	cli.getJob(2)
}
