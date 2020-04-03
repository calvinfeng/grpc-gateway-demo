package robotallocator

import (
	"context"
	"errors"
	"github.com/calvinfeng/grpc-gateway-demo/protos/robotrpc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
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

func (a *allocator) LeaseRobot(ctx context.Context, req *robotrpc.RobotLeaseRequest) (*robotrpc.RobotLeaseGrant, error) {
	time.Sleep(5 * time.Second)
	deadline, done := ctx.Deadline()
	if done {
		logrus.Warnf("deadline %s has reached", deadline)
	}

	select {
	case <-ctx.Done():
		logrus.Warnf("confirmed that context is dead")
		default:
			logrus.Info("context is healthy")
	}
	return nil, errors.New("context has reached deadline")
}

func (a *allocator) ListRobots(ctx context.Context, q *robotrpc.RobotQuery) (*robotrpc.RobotQueryResult, error) {
	res := new(robotrpc.RobotQueryResult)

	if q.RobotType != "Terminator" {
		return nil, status.Errorf(codes.OutOfRange, "%s is not a valid robot type", q.RobotType)
	}

	res.Robots = a.robots
	return res, nil
}

