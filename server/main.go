package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const (
	PORT = ":55555" //may need this in shared folder with client
	user = "user"
	pass = "user"
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

type Post struct {
	Name string
	Args string
}

type GetResponse struct {
	Id  int
	Log string
}

func commandHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p Post
	err := decoder.Decode(&p)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case "POST":
		path, err := GetJobPath(p.Name)
		if err != nil {
			http.NotFound(w, r)
			log.Println("Unable to find job: ", p.Name)
			return
		}

		id := counter.Increment()
		job := &Job{id: id, log: ""}
		runLog[id] = job

		job.Start(path, p.Args)
		fmt.Fprintf(w, "{\"id\":%d}", id)
		log.Printf("DEBUG::: runLog - %v", runLog)
		return
	case "GET":
		log.Printf("GET request in /command/ ")

		to_send := make([]string, len(runLog))
		i := 0
		for k, v := range runLog {
			to_send[i] = fmt.Sprintf("{\"id\":%d,\"log\":\"%v\"}", k, v)
			i++
		}
		fmt.Fprintf(w, "[%v]", strings.Join(to_send, ","))
		return
	}
	http.NotFound(w, r)
}

func idHandler(w http.ResponseWriter, r *http.Request) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		log.Println("Request not found: ", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	id, err := strconv.Atoi(m[1])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	job, ok := runLog[id]
	if !ok {
		http.NotFound(w, r)
		return
	}

	//May come back and encode as json, so things like quotes work in log
	switch r.Method {
	case "GET":
		getResponse := GetResponse{id, job.log}
		data, err := json.Marshal(getResponse)
		if err != nil {
			log.Printf("Unable to encode job.log as json %v", getResponse)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	case "DELETE":
		fmt.Fprintf(w, "{\"success\":%t}", job.Kill())
		return
	case "HEAD":
		fmt.Fprintf(w, "{\"status\":%v}", job.Status())
		return
	}
	http.NotFound(w, r)
}

func routeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request from %v: %v", r.Header.Get("X-Forwarded-For"), r.URL.Path)
	u, p, ok := r.BasicAuth()
	if !ok {
		log.Printf("Not using basic auth.")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if u != user || p != pass {
		log.Printf("Invalid credentials for %v", u)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if r.URL.Path == "/command/" {
		commandHandler(w, r)
	} else {
		idHandler(w, r)
	}
}

func main() {
	http.HandleFunc("/command/", routeHandler)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
