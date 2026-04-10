// Package main implements a client for Prober service.
package main

import (
	"context"
	"time"
	"flag"
	"log"

	pb "github.com/CodeYourFuture/immersive-go-course/grpc-client-server/prober"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr     = flag.String("addr", "localhost:50051", "the address to connect to")
	rep      = flag.Int("rep", 1, "number of repetitions")
	endpoint = flag.String("endpoint", "http://www.google.com", "the endpoint to prob, defaults to http://www.google.com")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProberClient(conn)

	// Contact the server and print out its response.
	// TODO: add a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// TODO: endpoint should be a flag
	// TODO: add number of times to probe
	r, err := c.DoProbes(ctx, &pb.ProbeRequest{Endpoint: *endpoint, RequestCount: int32(*rep)})

	if err != nil {
		log.Fatalf("could not probe: %v", err)
	}

	if r.GetErrorState() {
		log.Fatalf("could not probe: %s", r.GetErrorMessage())
	}
	log.Printf("Response Time: %f", r.GetAverageRespMsecs())
}
