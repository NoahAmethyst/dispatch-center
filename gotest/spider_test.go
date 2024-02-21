package gotest

import (
	"context"
	"github.com/NoahAmethyst/dispatch-center/cluster/rpc"
	"github.com/NoahAmethyst/dispatch-center/cluster/spider_svc"
	"github.com/NoahAmethyst/dispatch-center/proto/pb/spider_pb"
	"testing"
)

func Test_Odaily(t *testing.T) {
	// Initialize grpc client
	for _, rpcCli := range rpc.RpcCliList {
		rpc.InitGrpcCli(rpcCli)
	}
	cli := spider_svc.SvcCli()
	resp, err := cli.OdailyFeeds(context.Background(), &spider_pb.SpiderReq{})
	if err != nil {
		t.Error(err.Error())
	} else {
		t.Logf("odaily size:%d", len(resp.OdailyFeeds))
	}
}
