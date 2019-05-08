package main

import (
	"testing"

	"github.com/brotherlogic/keystore/client"
	"golang.org/x/net/context"
)

func InitTestServer() *Server {
	s := Init()
	s.SkipLog = true
	s.GoServer.KSclient = *keystoreclient.GetTestClient(".test")
	return s
}

func TestRunEmptyQueue(t *testing.T) {
	s := InitTestServer()

	err := s.runQueue(context.Background())
	if err != nil {
		t.Errorf("Unable to run empty queue: %v", err)
	}
}
