package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/brotherlogic/goserver/utils"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/executor/proto"

	//Needed to pull in gzip encoding init
	_ "google.golang.org/grpc/encoding/gzip"
)

func main() {
	entries, err := utils.ResolveAll("executor")
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

	for _, entry := range entries {
		if len(*ondeck) == 0 || entry.Identifier == *ondeck {
			conn, err := grpc.Dial(entry.Ip+":"+strconv.Itoa(int(entry.Port)), grpc.WithInsecure())
			defer conn.Close()

			if err != nil {
				log.Fatalf("Unable to dial: %v", err)
			}

			client := pb.NewExecutorServiceClient(conn)

			resp, err := client.QueueExecute(ctx, &pb.ExecuteRequest{Command: &pb.Command{Binary: os.Args[adjust], Parameters: os.Args[adjust+1:]}})
			if err != nil {
				fmt.Printf("%v failed: %v\n", entry.Identifier, err)
			} else {
				for resp.Status != pb.CommandStatus_COMPLETE {
					fmt.Printf("%v %v\n", entry.Identifier, resp)
					time.Sleep(time.Second)
				}
			}
		}
	}
}
