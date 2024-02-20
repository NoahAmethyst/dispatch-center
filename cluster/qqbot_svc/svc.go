package qqbot_svc

import (
	"github.com/NoahAmethyst/dispatch-center/cluster/rpc"
	"github.com/NoahAmethyst/dispatch-center/proto/pb/qqbot_pb"
)

func SvcCli() qqbot_pb.QQBotServiceClient {
	return qqbot_pb.NewQQBotServiceClient(rpc.GetConn(rpc.CliQQBot))
}
