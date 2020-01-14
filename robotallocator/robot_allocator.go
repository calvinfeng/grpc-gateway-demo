package robotallocator

import (
	"context"
	"github.com/calvinfeng/grpc-gateway-demo/protos/robotrpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RobotAllocator interface {
	robotrpc.RobotAllocationServer
}

func New() RobotAllocator {
	return &allocator{robots: []string{"Robot 1", "Robot 2", "Robot 3"}}
}

type allocator struct {
	robots []string
}

func (a *allocator) LeaseRobot(context.Context, *robotrpc.RobotLeaseRequest) (*robotrpc.RobotLeaseGrant, error) {
	panic("implement me")
}

func (a *allocator) ListRobots(ctx context.Context, q *robotrpc.RobotQuery) (*robotrpc.RobotQueryResult, error) {
	res := new(robotrpc.RobotQueryResult)

	if q.RobotType != "Terminator" {
		return nil, status.Errorf(codes.OutOfRange, "%s is not a valid robot type", q.RobotType)
	}

	res.Robots = a.robots
	return res, nil
}

