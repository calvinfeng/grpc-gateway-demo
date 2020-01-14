package main

import (
	"context"
	"fmt"
	"github.com/calvinfeng/grpc-gateway-demo/protos/robotrpc"
	"github.com/calvinfeng/grpc-gateway-demo/robotallocator"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
)

const gRPCEndpoint = "localhost:9000"
const gatewayHTTPPort = ":8080"
const grpcPort = ":9000"

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:             true,
	})
}

func runGatewayServer() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	endpoint := fmt.Sprintf("localhost%s", grpcPort)
	err := robotrpc.RegisterRobotAllocationHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("launching HTTP server")
	if err := http.ListenAndServe(gatewayHTTPPort, mux); err != nil {
		logrus.Fatal(err)
	}
}

func runGRPCServer() {
	srv := grpc.NewServer()
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		logrus.Fatal(err)
	}

	// Register services
	robotrpc.RegisterRobotAllocationServer(srv, robotallocator.New())

	logrus.Info("launching gRPC server")
	if err := srv.Serve(lis); err != nil {
		logrus.Fatal(err)
	}
}

func main() {
	go runGRPCServer()
	runGatewayServer()
}