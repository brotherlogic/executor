package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/brotherlogic/goserver/utils"
	"google.golang.org/grpc"

	pbd "github.com/brotherlogic/discovery/proto"
	pb "github.com/brotherlogic/executor/proto"

	//Needed to pull in gzip encoding init
	_ "google.golang.org/grpc/encoding/gzip"
)

func run(ctx context.Context, client pb.ExecutorServiceClient, binary string, params []string, entry *pbd.RegistryEntry, c chan bool) {
	var err error
	var resp *pb.ExecuteResponse
	currState := pb.CommandStatus_COMPLETE
	for resp == nil || resp.Status != pb.CommandStatus_COMPLETE {
		resp, err = client.QueueExecute(ctx, &pb.ExecuteRequest{Command: &pb.Command{Binary: binary, Parameters: params}})
		if err != nil {
			fmt.Printf("%v failed: %v\n", entry.Identifier, err)
		} else {
			if resp.Status != currState {
				fmt.Printf("%v %v\n", entry.Identifier, resp)
				currState = resp.Status
			}
			time.Sleep(time.Second)

		}
	}
	fmt.Printf("DONE %v %v\n", entry.Identifier, resp)
	c <- true

}

func main() {
	entries, err := utils.BaseResolveAll("executor")
	if err != nil {
		log.Fatalf("Unable to reach organiser: %v", err)
	}
	ctx, cancel := utils.BuildContext("executor-cli", "executor")
	defer cancel()

	var ondeck = flag.String("server", "", "The server to run on")
	flag.Parse()

	adjust := 1
	if len(*ondeck) > 0 {
		if os.Args[1] == "--server" {
			adjust += 2
		}
	}

	for i, v := range os.Args {
		if v == "--server" && i != 1 {
			log.Fatalf("Flag must appear first")
		}
	}

	c := make(chan bool)
	count := 0
	for _, entry := range entries {
		if len(*ondeck) == 0 || entry.Identifier == *ondeck {
			conn, err := grpc.Dial(entry.Ip+":"+strconv.Itoa(int(entry.Port)), grpc.WithInsecure())
			defer conn.Close()

			if err != nil {
				log.Fatalf("Unable to dial: %v", err)
			}

			client := pb.NewExecutorServiceClient(conn)
			go run(ctx, client, os.Args[adjust], os.Args[adjust+1:], entry, c)
			count++
		}
	}

	for count > 0 {
		<-c
		count--
	}
}
