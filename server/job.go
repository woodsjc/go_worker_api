package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

type Job struct {
	id         int
	log        string
	process_id int
}

func (j Job) Kill() bool {
	if j.process_id == 0 {
		return false
	}

	err := syscall.Kill(j.process_id, syscall.SIGINT)
	if err != nil {
		return false
	}
	return true
}

func (j Job) Status() string {
	if j.process_id == 0 {
		return "completed"
	}

	p, err := os.FindProcess(j.process_id)
	if err != nil {
		return "completed"
	}

	err = p.Signal(syscall.Signal(0))
	if err != nil {
		return "completed"
	}
	return "running"
}

//Should this be run as goroutine
func (j *Job) Start(path string, args string) {
	cmd := exec.Command(path, args)

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Unable to get stdin/out for job %d.", j.id)
		return
	}

	//will need lock around this if call anywhere else
	j.process_id = cmd.Process.Pid
	err = cmd.Run()
	if err == nil {
		j.log = fmt.Sprintf("{%v}", err)
	} else {
		j.log = string(out)
	}
	j.process_id = 0

	log.Printf("DEBUG:::%v %v finished with output: %v", path, args, j.log)
}

func GetJobPath(job string) (string, error) {
	path, err := exec.LookPath(job)
	if err != nil {
		return "", err
	}

	return path, nil
}
