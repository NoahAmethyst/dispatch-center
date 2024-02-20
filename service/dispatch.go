package service

import (
	"context"
	"github.com/NoahAmethyst/dispatch-center/proto/pb/dispatch_pb"
)

type DispatchServer struct{}

func (s DispatchServer) Push(ctx context.Context, req *dispatch_pb.Message) (resp *dispatch_pb.Resp, err error) {
	
	return
}
