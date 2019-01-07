package main

import (
	"log"
	"os"
	"sync"
	"testing"

	pb "github.com/brotherlogic/executor/proto"
)

func dlog(str string) {
	log.Printf("%v", str)
}

func TestSchedulerRun(t *testing.T) {
	os.Unsetenv("GOBIN")
	os.Unsetenv("GOPATH")

	s := Scheduler{
		commands:     make([]*rCommand, 0),
		executeMutex: &sync.Mutex{},
		log:          dlog,
	}

	output, err := s.schedule(&pb.Command{Binary: "ls", Parameters: []string{"-ltr"}})

	if err != nil {
		t.Errorf("Error running ls command: %v", err)
	}

	if len(output) == 0 {
		t.Errorf("No output produced")
	}

}

func TestBadSchedulerRun(t *testing.T) {
	s := Scheduler{
		commands:     make([]*rCommand, 0),
		executeMutex: &sync.Mutex{},
		log:          dlog,
	}

	output, err := s.schedule(&pb.Command{Binary: "madeupcommand", Parameters: []string{"-ltr"}})

	if err == nil {
		t.Errorf("No error running comand: %v", output)
	}
}

func TestStdErrSchedulerRu(t *testing.T) {
	s := Scheduler{
		commands:     make([]*rCommand, 0),
		executeMutex: &sync.Mutex{},
		log:          dlog,
	}

	_, err := s.schedule(&pb.Command{Binary: "./run.sh", Parameters: []string{}})

	if err != nil {
		t.Errorf("Unable to run simple err command: %v", err)
	}

}
