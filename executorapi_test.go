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
