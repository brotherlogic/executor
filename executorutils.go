package main

import (
	"fmt"
	"time"

	pb "github.com/brotherlogic/executor/proto"
	"github.com/brotherlogic/goserver/utils"
	"google.golang.org/grpc/status"
)

func (s *Server) runQueue() {
	for entry := range s.queue {
		ctx, cancel := utils.ManualContext("queue-run", time.Minute)
		Backlog.Set(float64(len(s.queue)))
		s.CtxLog(ctx, fmt.Sprintf("THE QUEUE EXEC OUT START: %+v => %v", entry, entry.req.Command.Binary))
		entry.resp.Status = pb.CommandStatus_IN_PROGRESS
		output, err := s.runExecute(ctx, entry.req)
		s.CtxLog(ctx, fmt.Sprintf("THE QUEUE EXEC OUT COMPLETE: %+v => %v (%v)", entry, entry.req.Command.Binary, err))
		if err != nil {
			entry.resp.CommandOutput = fmt.Sprintf("%v", err)
			entry.resp.ExitCode = int32(status.Convert(err).Code())
		} else {
			entry.resp.CommandOutput = output
		}
		entry.resp.Status = pb.CommandStatus_COMPLETE

		s.CtxLog(ctx, "Acking Channel")
		entry.ack <- true
		s.CtxLog(ctx, "Acked Channel")
		cancel()

	}

	s.done <- true
}
