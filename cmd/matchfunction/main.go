package main

import (
	"log"
	"net"
	"os"

	openmatch_mmf_testing "github.com/castaneai/openmatch-mmf-testing"

	"google.golang.org/grpc"
	"open-match.dev/open-match/pkg/pb"
)

func main() {
	qsAddr := os.Getenv("OPENMATCH_QUERY_SERVICE_ADDR")
	qsc, err := newQueryServiceClient(qsAddr)
	if err != nil {
		log.Fatalf("failed to connect to QueryService: %+v", err)
	}

	addr := ":50502"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %+v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMatchFunctionServer(s, openmatch_mmf_testing.NewMatchFunctionService(qsc))

	log.Printf("litening on %s...", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %+v", err)
	}
}

func newQueryServiceClient(addr string) (pb.QueryServiceClient, error) {
	cc, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return pb.NewQueryServiceClient(cc), nil
}
