package main

import (
	"testing"

	"github.com/brotherlogic/keystore/client"
)

func InitTestServer() *Server {
	s := Init()
	s.SkipLog = true
	s.GoServer.KSclient = *keystoreclient.GetTestClient(".test")
	return s
}

func TestBlank(t *testing.T) {
	blank()
}
