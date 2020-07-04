package main

import (
	"github.com/brotherlogic/keystore/client"
)

func InitTestServer() *Server {
	s := Init()
	s.SkipLog = true
	s.GoServer.KSclient = *keystoreclient.GetTestClient(".test")

	go func() {
		s.runQueue()
	}()

	return s
}
