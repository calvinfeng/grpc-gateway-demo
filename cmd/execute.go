package cmd

import (
	"context"
	"fmt"
	"github.com/calvinfeng/grpc-gateway-demo/protos/robotrpc"
	"github.com/calvinfeng/grpc-gateway-demo/robotallocator"
	gatewayRuntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"runtime"
	"strings"
	"time"
)

func callerPrettify(f *runtime.Frame) (function string, file string) {
	parts := strings.Split(f.Function, ".")
	fnName := parts[len(parts)-1]
	return fmt.Sprintf("%s:", fnName), fmt.Sprintf("[%s:%d]", f.File, f.Line)
}

func configureLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:               true,
		FullTimestamp:             true,
		CallerPrettyfier:          callerPrettify,
	})
	logrus.SetReportCaller(true)
}

func Execute() error {
	cobra.OnInitialize(configureLogger)
	root := &cobra.Command{
		Use: "grpc-gateway-demo",
		Short: "run a demo on how to setup gRPC server for HTTP requests",
	}

	root.AddCommand(
		&cobra.Command{
			Use: "runserver",
			Short: "run gRPC server",
			RunE: runServer,
		},
		&cobra.Command{
				Use: "runclient",
				Short: "run a client to hit gRPC server",
				RunE: runClient,
		},
		)

	return root.Execute()
}

const gatewayHTTPPort = ":8080"
const gRPCPort = ":9000"

func runGatewayServer() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := gatewayRuntime.NewServeMux()
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

func runServer(_ *cobra.Command, _ []string) error {
	go runGRPCServer()
	runGatewayServer()
	return nil
}

func runClient(_ *cobra.Command, _ []string) error {
	conn, err := grpc.Dial("localhost" + gRPCPort, grpc.WithInsecure())
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}

	cli := robotrpc.NewRobotAllocationClient(conn)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
	grant, err := cli.LeaseRobot(ctx, &robotrpc.RobotLeaseRequest{
		RobotNameId:          "freight100-001",
	})
	defer cancel()

	if err != nil {
		logrus.WithError(err).Fatal("robot lease failed")
	}

	logrus.Infof("robot is leased id=%d", grant.LeaseId)
	return nil
}