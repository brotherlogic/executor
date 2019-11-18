package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/brotherlogic/goserver"
	"github.com/brotherlogic/keystore/client"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/executor/proto"
	pbg "github.com/brotherlogic/goserver/proto"
)

type queueEntry struct {
	req  *pb.ExecuteRequest
	resp *pb.ExecuteResponse
}

//Server main server type
type Server struct {
	*goserver.GoServer
	scheduler *Scheduler
	queue     []*queueEntry
}

// Init builds the server
func Init() *Server {
	s := &Server{
		&goserver.GoServer{},
		&Scheduler{
			commands:     make([]*rCommand, 0),
			executeMutex: &sync.Mutex{},
		},
		make([]*queueEntry, 0),
	}
	s.scheduler.log = s.Log
	return s
}

// DoRegister does RPC registration
func (s *Server) DoRegister(server *grpc.Server) {
	pb.RegisterExecutorServiceServer(server, s)
}

// ReportHealth alerts if we're not healthy
func (s *Server) ReportHealth() bool {
	return true
}

// Shutdown the server
func (s *Server) Shutdown(ctx context.Context) error {
	return nil
}

// Mote promotes/demotes this server
func (s *Server) Mote(ctx context.Context, master bool) error {
	return nil
}

// GetState gets the state of the server
func (s *Server) GetState() []*pbg.State {
	v := []string{}
	for _, q := range s.queue {
		v = append(v, fmt.Sprintf("%v", q.resp.Status))
	}
	return []*pbg.State{
		&pbg.State{Key: "queue_size", Value: int64(len(s.queue))},
		&pbg.State{Key: "runs", Value: s.scheduler.runs},
		&pbg.State{Key: "state", Text: fmt.Sprintf("%v", v)},
	}
}

func main() {
	var quiet = flag.Bool("quiet", false, "Show all output")
	flag.Parse()

	//Turn off logging
	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	server := Init()
	server.PrepServer()
	server.Register = server

	server.GoServer.KSclient = *keystoreclient.GetClient(server.DialMaster)

	server.RegisterServerIgnore("executor", false, true)

	server.RegisterRepeatingTaskNonMaster(server.runQueue, "run_queue", time.Minute)

	fmt.Printf("%v", server.Serve())
}
