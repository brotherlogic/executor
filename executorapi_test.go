package main

import (
	"context"
	"testing"

	pb "github.com/brotherlogic/executor/proto"
)

func TestRunAPI(t *testing.T) {
	s := InitTestServer()
	resp, err := s.Execute(context.Background(), &pb.ExecuteRequest{Command: &pb.Command{Binary: "ls"}})

	if err != nil {
		t.Fatalf("Error running command: %v", err)
	}

	if resp.TimeTakenInMillis == 0 {
		t.Errorf("Command took no time")
	}
}

func TestRunQueueAPI(t *testing.T) {
	s := InitTestServer()
	_, err := s.QueueExecute(context.Background(), &pb.ExecuteRequest{Command: &pb.Command{Binary: "ls", Parameters: []string{"-ltr"}}})

	if err != nil {
		t.Errorf("Failed: %v", err)
	}

	resp, err := s.QueueExecute(context.Background(), &pb.ExecuteRequest{Command: &pb.Command{Binary: "ls", Parameters: []string{"-ltr"}}})

	if err != nil {
		t.Errorf("Failed: %v", err)
	}

	s.drainQueue()

	if resp.Status != pb.CommandStatus_COMPLETE {
		t.Errorf("Bad resp: %v", resp)
	}
}

func TestRunQueueWithBadCommand(t *testing.T) {
	s := InitTestServer()
	_, err := s.QueueExecute(context.Background(), &pb.ExecuteRequest{Command: &pb.Command{Binary: "ltttts", Parameters: []string{"-ltr"}}})

	if err != nil {
		t.Errorf("Failed: %v", err)
	}

	resp, err := s.QueueExecute(context.Background(), &pb.ExecuteRequest{Command: &pb.Command{Binary: "ltttts", Parameters: []string{"-ltr"}}})

	if err != nil {
		t.Errorf("Failed: %v", err)
	}

	s.drainQueue()

	if resp.Status != pb.CommandStatus_COMPLETE {
		t.Errorf("Bad resp: %v", resp)
	}
}

func TestMini(t *testing.T) {
	if mini(10, 5) != 5 {
		t.Errorf("Bad min")
	}
}
