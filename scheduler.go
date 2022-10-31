package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	pb "github.com/brotherlogic/executor/proto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/net/context"
)

var (
	execLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "executor_latency",
		Help:    "The latency of server requests",
		Buckets: []float64{0.5, 1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024},
	}, []string{"key"})
)

// Scheduler for doing builds
type Scheduler struct {
	commands     []*rCommand
	runs         int64
	executeMutex *sync.Mutex
	log          func(ctx context.Context, str string)
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

func (s *Scheduler) schedule(ctx context.Context, command *pb.Command, key string) (string, error) {
	s.executeMutex.Lock()
	defer s.executeMutex.Unlock()

	rCommand := &rCommand{
		base:    command,
		command: exec.Command(command.Binary, command.Parameters...),
	}

	s.log(ctx, fmt.Sprintf("Running command: %v", command.Binary))
	t1 := time.Now()
	s.runAndWait(rCommand)
	execLatency.With(prometheus.Labels{"key": key}).Observe(float64(time.Since(t1).Seconds()))
	s.log(ctx, fmt.Sprintf("%v took %v", command.Binary, time.Since(t1)))

	s.log(ctx, fmt.Sprintf("Ran: %v, %v -> %v %v", command.Binary, command.Parameters, rCommand.output, rCommand.err))
	return rCommand.output, rCommand.err
}

func (s *Scheduler) runAndWait(c *rCommand) {
	c.err = s.run(c, true)
}

func (s *Scheduler) run(c *rCommand, hardwait bool) error {
	s.runs++

	// Setup the gopath
	env := os.Environ()
	gpath := "/home/simon/code"
	c.command.Path = strings.Replace(c.command.Path, "$GOPATH", gpath, -1)
	for i := range c.command.Args {
		c.command.Args[i] = strings.Replace(c.command.Args[i], "$GOPATH", gpath, -1)
	}
	path := fmt.Sprintf("GOPATH=/home/simon/code")
	found := false
	for i, blah := range env {
		if strings.HasPrefix(blah, "GOPATH") {
			env[i] = path
			found = true
		}
	}
	if !found {
		env = append(env, path)
	}
	c.command.Env = env

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
	runner := func() {
		err := c.command.Wait()
		c.endTime = time.Now().Unix()

		if err != nil {
			c.err = fmt.Errorf("%v -> %v / %v", err, c.output, c.erroutput)
		}
	}
	if !hardwait {
		go runner()
	} else {
		runner()
	}

	return nil
}
