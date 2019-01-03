package main

import (
	"bufio"
	"os/exec"
	"sync"
	"time"

	pb "github.com/brotherlogic/executor/proto"
)

// Scheduler for doing builds
type Scheduler struct {
	commands     []*rCommand
	runs         int64
	executeMutex *sync.Mutex
}

type rCommand struct {
	base      *pb.Command
	command   *exec.Cmd
	output    string
	erroutput string
	startTime int64
	endTime   int64
	err       error
}

func (s *Scheduler) schedule(command *pb.Command) (string, error) {
	s.executeMutex.Lock()
	defer s.executeMutex.Unlock()

	rCommand := &rCommand{
		base:    command,
		command: exec.Command(command.Binary, command.Parameters...),
	}

	s.runAndWait(rCommand)
	return rCommand.output, rCommand.err
}

func (s *Scheduler) runAndWait(c *rCommand) {
	c.err = s.run(c)
	if c.err == nil {
		for c.endTime == 0 {
			time.Sleep(time.Second)
		}
	}
}

func (s *Scheduler) run(c *rCommand) error {
	out, _ := c.command.StdoutPipe()
	if out != nil {
		scanner := bufio.NewScanner(out)
		go func() {
			for scanner != nil && scanner.Scan() {
				c.output += scanner.Text()
			}
			out.Close()
		}()
	}

	out2, _ := c.command.StderrPipe()
	if out2 != nil {
		scanner := bufio.NewScanner(out2)
		go func() {
			for scanner != nil && scanner.Scan() {
				c.erroutput += scanner.Text()
			}
			out2.Close()
		}()
	}

	err := c.command.Start()
	if err != nil {
		return err
	}
	c.startTime = time.Now().Unix()

	// Monitor the job and report completion
	go func() {
		err := c.command.Wait()
		c.endTime = time.Now().Unix()
		c.err = err
	}()

	return nil
}
