package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	pb "github.com/CodeYourFuture/immersive-go-course/grpc-client-server/prober"
	"google.golang.org/grpc"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	// "github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

var latencyGuage = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "latency_gauge_value",
	Help: "This gauge the latency of an endpoint",
}, []string{"endpoint"})

// server is used to implement prober.ProberServer.
type server struct {
	pb.UnimplementedProberServer
}

func (s *server) DoProbes(ctx context.Context, in *pb.ProbeRequest) (*pb.ProbeReply, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// TODO: support a number of repetitions and return average latency
	var wg sync.WaitGroup
	var mu sync.Mutex
	var totalElapsedTime float32
	var errState bool
	var errMessage string
	requestCount := int(in.GetRequestCount())

	for i := 0; i < requestCount; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			start := time.Now()

			_, err := http.Get(in.GetEndpoint())
			if err != nil {
				mu.Lock()
				errState = true
				errMessage = err.Error()
				mu.Unlock()
				return
			}

			elapsed := time.Since(start)
			elapsedMsecs := float32(elapsed / time.Millisecond)

			mu.Lock()
			totalElapsedTime += elapsedMsecs
			mu.Unlock()
		}()
	}

	//    start := time.Now()
	// _, _ = http.Get(in.GetEndpoint())	// TODO: add error handling here and check the response code
	// elapsed := time.Since(start)
	// elapsedMsecs := float32(elapsed / time.Millisecond)
	wg.Wait()

	if errState {
		return &pb.ProbeReply{ErrorState: true, ErrorMessage: errMessage}, nil
	}
	averageRespMsecs := totalElapsedTime / float32(in.GetRequestCount())

	latencyGuage.WithLabelValues(in.GetEndpoint()).Set(float64(averageRespMsecs))

	return &pb.ProbeReply{AverageRespMsecs: averageRespMsecs}, nil
}

func initiateGrpcServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterProberServer(s, &server{})
	log.Printf("grpc server listening on %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func initiateMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}

func main() {

	flag.Parse()
	go initiateGrpcServer()
	go initiateMetricsServer()

	// Block forever
	select {}

	// lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	// if err != nil {
	// 	log.Fatalf("failed to listen: %v", err)
	// }
	// s := grpc.NewServer()
	// pb.RegisterProberServer(s, &server{})
	// log.Printf("server listening at %v", lis.Addr())
	// if err := s.Serve(lis); err != nil {
	// 	log.Fatalf("failed to serve: %v", err)
	// }
	//
}
