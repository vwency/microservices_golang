package transport

import (
	"context"

	"github.com/go-kit/kit/log"
	gokitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"

	"github.com/vwency/microservices_golang/internal/hello_service/endpoint"
	pb "github.com/vwency/microservices_golang/proto/hello_service"
)

type grpcServer struct {
	sayHello gokitgrpc.Handler
	pb.UnimplementedHelloServiceServer
}

func NewGRPCServer(endpoints endpoint.Endpoints, logger log.Logger) pb.HelloServiceServer {
	options := []gokitgrpc.ServerOption{
		gokitgrpc.ServerErrorLogger(logger),
	}

	return &grpcServer{
		sayHello: gokitgrpc.NewServer(
			endpoints.SayHelloEndpoint,
			decodeSayHelloRequest,
			encodeSayHelloResponse,
			options...,
		),
	}
}

func (s *grpcServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	_, rep, err := s.sayHello.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.HelloResponse), nil
}

func decodeSayHelloRequest(_ context.Context, req interface{}) (interface{}, error) {
	r := req.(*pb.HelloRequest)
	return endpoint.SayHelloRequest{Text: r.GetText()}, nil
}

func encodeSayHelloResponse(_ context.Context, resp interface{}) (interface{}, error) {
	r := resp.(endpoint.SayHelloResponse)
	return &pb.HelloResponse{Text: r.Text}, r.Err
}

func RegisterHelloServiceServer(grpcServer *grpc.Server, server pb.HelloServiceServer) {
	pb.RegisterHelloServiceServer(grpcServer, server)
}
