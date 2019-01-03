package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/brotherlogic/goserver/utils"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/executor/proto"
	pbgs "github.com/brotherlogic/goserver/proto"
	pbt "github.com/brotherlogic/tracer/proto"

	//Needed to pull in gzip encoding init
	_ "google.golang.org/grpc/encoding/gzip"
)

func main() {
	entries, err := utils.ResolveAll("executor")
	if err != nil {
		log.Fatalf("Unable to reach organiser: %v", err)
	}
	ctx, cancel := utils.BuildContext("recordwants-cli", "recordwants", pbgs.ContextType_LONG)
	defer cancel()

	for _, entry := range entries {
		conn, err := grpc.Dial(entry.Ip+":"+strconv.Itoa(int(entry.Port)), grpc.WithInsecure())
		defer conn.Close()

		if err != nil {
			log.Fatalf("Unable to dial: %v", err)
		}

		client := pb.NewExecutorServiceClient(conn)
		resp, err := client.Execute(ctx, &pb.ExecuteRequest{Command: &pb.Command{Binary: os.Args[1], Parameters: os.Args[2:]}})
		if err != nil {
			fmt.Printf("%v failed: %v\n", entry.Identifier, err)
		} else {
			fmt.Printf("%v (%v): %v\n", entry.Identifier, resp.TimeTakenInMillis, resp.CommandOutput)
		}
	}
	utils.SendTrace(ctx, "End of CLI", time.Now(), pbt.Milestone_END, "recordwants-cli")
}
