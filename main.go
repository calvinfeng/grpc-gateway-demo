package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/calvinfeng/grpc-gateway-demo/protos/robotrpc"
	"github.com/calvinfeng/grpc-gateway-demo/robotallocator"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const gatewayHTTPPort = ":8080"
const gRPCPort = ":9000"

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func runGatewayServer() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	endpoint := fmt.Sprintf("localhost%s", gRPCPort)
	err := robotrpc.RegisterRobotAllocationHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Infof("launching gRPC Gateway HTTP server on %s", gatewayHTTPPort)
	if err := http.ListenAndServe(gatewayHTTPPort, mux); err != nil {
		logrus.Fatal(err)
	}
}

func runGRPCServer() {
	srv := grpc.NewServer()
	lis, err := net.Listen("tcp", gRPCPort)
	if err != nil {
		logrus.Fatal(err)
	}

	// Register services
	robotrpc.RegisterRobotAllocationServer(srv, robotallocator.New())

	logrus.Infof("launching gRPC server on %s", gRPCPort)
	if err := srv.Serve(lis); err != nil {
		logrus.Fatal(err)
	}
}

func main() {
	go runGRPCServer()
	runGatewayServer()
}
