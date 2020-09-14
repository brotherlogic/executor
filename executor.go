package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/brotherlogic/goserver"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/executor/proto"
	pbg "github.com/brotherlogic/goserver/proto"
)

type queueEntry struct {
	req  *pb.ExecuteRequest
	resp *pb.ExecuteResponse
	ack  chan bool
}

var (
	//Backlog - the print queue
	Backlog = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "executor_backlog",
		Help: "The size of the executor queue",
	})

	archive = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "executor_archive",
		Help: "The size of the executor queue",
	})
)

//Server main server type
type Server struct {
	*goserver.GoServer
	scheduler *Scheduler
	queue     chan *queueEntry
	archive   []*queueEntry
	done      chan bool
}

// Init builds the server
func Init() *Server {
	s := &Server{
		&goserver.GoServer{},
		&Scheduler{
			commands:     make([]*rCommand, 0),
			executeMutex: &sync.Mutex{},
		},
		make(chan *queueEntry, 10000),
		make([]*queueEntry, 0),
		make(chan bool),
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
	return []*pbg.State{
		&pbg.State{Key: "nope", Value: int64(123)},
	}
}

func (s *Server) drainQueue() {
	close(s.queue)
	<-s.done
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

	err := server.RegisterServerV2("executor", false, true)
	if err != nil {
		return
	}

	go func() {
		server.runQueue()
	}()

	fmt.Printf("%v", server.Serve())
}
