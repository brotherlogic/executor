package main

import (
	"fmt"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/executor/proto"
)

func (s *Server) runQueue(ctx context.Context) error {
	for _, entry := range s.queue {
		s.Log(fmt.Sprintf("Attemping %v", entry))
		if entry.resp.Status == pb.CommandStatus_IN_QUEUE {
			entry.resp.Status = pb.CommandStatus_IN_PROGRESS
			output, err := s.runExecute(ctx, entry.req)
			if err != nil {
				entry.resp.CommandOutput = fmt.Sprintf("%v", err)
			} else {
				entry.resp.CommandOutput = output
			}
			entry.resp.Status = pb.CommandStatus_COMPLETE
			return nil
		}
	}

	return nil
}
