package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/TylerJGabb/go-calc-grpc-contract/contract"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct {
	contract.UnimplementedCalculatorServer
}

func (s *server) Add(
	ctx context.Context,
	in *contract.CalculationRequest,
) (*contract.CalculationResponse, error) {
	fmt.Printf("ADD: %d+%d\n", in.A, in.B)
	return &contract.CalculationResponse{
		Result: in.A + in.B,
	}, nil
}

func (s *server) Divide(
	ctx context.Context,
	in *contract.CalculationRequest,
) (*contract.CalculationResponse, error) {
	fmt.Printf("DIVIDE: %d/%d\n", in.A, in.B)
	if in.B == 0 {
		fmt.Printf("DIVIDE: cannot divide by zero\n")
		return nil, status.Error(codes.InvalidArgument, "cannot divide by zero")
	}
	return &contract.CalculationResponse{
		Result: in.A / in.B,
	}, nil
}

func (s *server) Sum(
	ctx context.Context,
	in *contract.NumbersRequest,
) (*contract.CalculationResponse, error) {
	fmt.Printf("SUM: %v\n", in.Numbers)
	var sum int64
	for _, v := range in.Numbers {
		sum += v
	}
	return &contract.CalculationResponse{
		Result: sum,
	}, nil
}

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	contract.RegisterCalculatorServer(s, &server{})
	fmt.Println("Starting server on port " + port)
	if err := s.Serve(listener); err != nil {
		panic(err)
	}
}
