package main

import (
	"context"
	"github.com/NoahAmethyst/dispatch-center/cluster"
	"github.com/NoahAmethyst/dispatch-center/cluster/rpc"
	"github.com/NoahAmethyst/dispatch-center/constant"
	"github.com/NoahAmethyst/dispatch-center/proto/pb/dispatch_pb"
	"github.com/NoahAmethyst/dispatch-center/task"
	"github.com/NoahAmethyst/dispatch-center/utils/cron"
	"os"
	"time"
)

func main() {
	//Set time location to East eighth District
	time.Local = time.FixedZone("UTC", 8*60*60)

	ctx := context.Background()

	gracefulShutdown(ctx, cluster.KubeOptServer)

	// Initialize grpc client
	for _, rpcCli := range rpc.RpcCliList {
		rpc.InitGrpcCli(rpcCli)
	}

	task.RegisterTask(ctx, dispatch_pb.Bot_Dingding, cron.JobDurationMinutely)
	// Start Grpc server
	cluster.StartServer(os.Getenv(constant.GrpcListenPort))
}
