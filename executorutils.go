package main

import (
	"fmt"

	pb "github.com/brotherlogic/executor/proto"
)

func (s *Server) runQueue() {
	for entry := range s.queue {
		s.Log(fmt.Sprintf("THE QUEUE EXEC OUT START: %+v => %v", entry, entry.req.Command.Binary))
		entry.resp.Status = pb.CommandStatus_IN_PROGRESS
		output, err := s.runExecute(entry.req)
		s.Log(fmt.Sprintf("THE QUEUE EXEC OUT COMPLETE: %+v => %v", entry, entry.req.Command.Binary))
		if err != nil {
			entry.resp.CommandOutput = fmt.Sprintf("%v", err)
		} else {
			entry.resp.CommandOutput = output
		}
		entry.resp.Status = pb.CommandStatus_COMPLETE

		entry.ack <- true
	}

	s.done <- true
}
