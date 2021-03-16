package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	USER = "user"
	PASS = "user"
	URL  = "http://localhost:55555/command/"
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

func basicAuth(r *http.Request) {
	r.SetBasicAuth(USER, PASS)
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

	basicAuth(request)
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

	basicAuth(request)
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

func main() {
	cli := &Cli{&http.Client{}}
	cli.postJob("ls", "-ax")
	cli.postJob("env", "")

	cli.getJob(1)
	cli.getJob(2)
}
