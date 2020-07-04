package main

import (
	"time"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/executor/proto"
)

func (s *Server) runExecute(req *pb.ExecuteRequest) (string, error) {
	return s.scheduler.schedule(req.Command)
}

// Execute executes a command
func (s *Server) Execute(ctx context.Context, req *pb.ExecuteRequest) (*pb.ExecuteResponse, error) {
	sTime := time.Now()
	output, err := s.scheduler.schedule(req.Command)

	return &pb.ExecuteResponse{
		TimeTakenInMillis: time.Now().Sub(sTime).Nanoseconds() / 100000,
		CommandOutput:     output,
	}, err
}

func mini(a, b int) int {
	if a > b {
		return b
	}
	return a
}

// QueueExecute executes a command
func (s *Server) QueueExecute(ctx context.Context, req *pb.ExecuteRequest) (*pb.ExecuteResponse, error) {
	// Pre clean the queue
	nq := []*queueEntry{}
	for _, q := range s.archive {
		if !q.req.GetReadyForDeletion() {
			nq = append(nq, q)
		}
	}
	s.archive = nq

	Backlog.Set(float64(len(s.queue)))

	for _, q := range s.archive {
		match := q.req.Command.Binary == req.Command.Binary && len(q.req.Command.Parameters) == len(req.Command.Parameters)
		for i := 0; i < mini(len(q.req.Command.Parameters), len(req.Command.Parameters)); i++ {
			match = match && q.req.Command.Parameters[i] == req.Command.Parameters[i]
		}

		if match {
			q.req.ReadyForDeletion = q.req.GetCommand().GetDeleteOnComplete()
			return q.resp, nil
		}
	}

	r := &pb.ExecuteResponse{Status: pb.CommandStatus_IN_QUEUE}
	entry := &queueEntry{req: req, resp: r, ack: make(chan bool, 100)}
	s.archive = append(s.archive, entry)
	s.queue <- entry
	return r, nil
}
